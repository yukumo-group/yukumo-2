package aq

import (
	"context"
	"fmt"
	"sync"
)

// Pool hands out Runtime instances for exclusive use by one goroutine at a time.
// Acquire a runtime, use it only on that goroutine, then Release it.
type Pool struct {
	mu     sync.Mutex
	idle   []*Runtime
	active int
	closed bool
	newFn  func(context.Context) (*Runtime, error)
}

func NewPool() *Pool {
	return &Pool{newFn: New}
}

// Acquire returns an exclusive Runtime. Caller must Release it.
func (p *Pool) Acquire(ctx context.Context) (*Runtime, error) {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil, fmt.Errorf("aq pool closed")
	}
	if n := len(p.idle); n > 0 {
		r := p.idle[n-1]
		p.idle = p.idle[:n-1]
		p.active++
		p.mu.Unlock()
		return r, nil
	}
	p.mu.Unlock()

	r, err := p.newFn(ctx)
	if err != nil {
		return nil, err
	}

	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		_ = r.Close(ctx)
		return nil, fmt.Errorf("aq pool closed")
	}
	p.active++
	p.mu.Unlock()
	return r, nil
}

// Release returns a Runtime to the pool. Do not use r after Release.
func (p *Pool) Release(r *Runtime) {
	if r == nil {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.active > 0 {
		p.active--
	}
	if p.closed {
		_ = r.Close(context.Background())
		return
	}
	p.idle = append(p.idle, r)
}

// Do acquires a runtime, runs fn, then releases it.
func (p *Pool) Do(ctx context.Context, fn func(*Runtime) error) error {
	r, err := p.Acquire(ctx)
	if err != nil {
		return err
	}
	defer p.Release(r)
	return fn(r)
}

func (p *Pool) HasPreset(ctx context.Context, version, preset int32) (bool, error) {
	var ok bool
	err := p.Do(ctx, func(r *Runtime) error {
		var err error
		ok, err = r.HasPreset(ctx, version, preset)
		return err
	})
	return ok, err
}

func (p *Pool) PresetNames(ctx context.Context, version int32) ([]string, error) {
	var names []string
	err := p.Do(ctx, func(r *Runtime) error {
		var err error
		names, err = r.PresetNames(ctx, version)
		return err
	})
	return names, err
}

func (p *Pool) Synthesize(ctx context.Context, text string, param Param) ([]byte, error) {
	var wav []byte
	err := p.Do(ctx, func(r *Runtime) error {
		var err error
		wav, err = r.Synthesize(ctx, text, param)
		return err
	})
	return wav, err
}

func (p *Pool) Close(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.closed = true
	var first error
	for _, r := range p.idle {
		if err := r.Close(ctx); err != nil && first == nil {
			first = err
		}
	}
	p.idle = nil
	return first
}
