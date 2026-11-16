package aq_test

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"testing"

	"yukumo-script/internal/domain"
	"yukumo-script/runtime/aq"
)

func TestSynthesizeWAV(t *testing.T) {
	ctx := context.Background()
	rt, err := aq.New(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = rt.Close(ctx) })

	wav, err := rt.Synthesize(ctx, "こんにちは", aq.Param{
		Version: domain.Version1,
		Preset:  0,
		Speed:   100,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(wav) < 44 || !bytes.HasPrefix(wav, []byte("RIFF")) {
		t.Fatalf("expected WAV, got %d bytes prefix=%q", len(wav), wav[:min(4, len(wav))])
	}
}

func TestPoolConcurrentSynthesize(t *testing.T) {
	ctx := context.Background()
	pool := aq.NewPool()
	t.Cleanup(func() { _ = pool.Close(ctx) })

	const n = 8
	var wg sync.WaitGroup
	errCh := make(chan error, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, err := pool.Acquire(ctx)
			if err != nil {
				errCh <- err
				return
			}
			defer pool.Release(r)
			wav, err := r.Synthesize(ctx, "あいうえお", aq.Param{
				Version: domain.Version1,
				Preset:  0,
				Speed:   100,
			})
			if err != nil {
				errCh <- err
				return
			}
			if len(wav) < 44 || !bytes.HasPrefix(wav, []byte("RIFF")) {
				errCh <- fmt.Errorf("invalid wav length %d", len(wav))
			}
		}()
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestPoolSynthesizeHelper(t *testing.T) {
	ctx := context.Background()
	pool := aq.NewPool()
	t.Cleanup(func() { _ = pool.Close(ctx) })

	wav, err := pool.Synthesize(ctx, "テスト", aq.Param{
		Version: domain.Version1,
		Preset:  0,
		Speed:   100,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(wav, []byte("RIFF")) {
		t.Fatalf("expected WAV prefix, got %q", wav[:min(4, len(wav))])
	}
}
