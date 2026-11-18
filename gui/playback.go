package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync"
)

// wavPlayer plays preview audio through an OS player instead of the WebView's
// media pipeline. WebKitGTK needs a GStreamer autoaudiosink for HTMLAudio,
// which is missing on many Linux desktops.
type wavPlayer struct {
	mu      sync.Mutex
	cmd     *exec.Cmd
	cleanup string
}

func newWAVPlayer() *wavPlayer {
	return &wavPlayer{}
}

func (p *wavPlayer) Play(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("preview: empty WAV data")
	}
	p.Stop()

	cmd, cleanup, err := startPlayer(data)
	if err != nil {
		return err
	}

	p.mu.Lock()
	p.cmd = cmd
	p.cleanup = cleanup
	p.mu.Unlock()

	go func() {
		_ = cmd.Wait()
		p.mu.Lock()
		if p.cmd == cmd {
			p.cmd = nil
			p.cleanup = ""
		}
		p.mu.Unlock()
		if cleanup != "" {
			_ = os.Remove(cleanup)
		}
	}()
	return nil
}

func (p *wavPlayer) Stop() {
	p.mu.Lock()
	cmd := p.cmd
	cleanup := p.cleanup
	p.cmd = nil
	p.cleanup = ""
	p.mu.Unlock()

	if cmd != nil && cmd.Process != nil {
		_ = cmd.Process.Kill()
	}
	if cleanup != "" {
		_ = os.Remove(cleanup)
	}
}

// startPlayer prefers stdin for Linux players so preview does not create a
// second file next to synthesized audio.
func startPlayer(data []byte) (*exec.Cmd, string, error) {
	if runtime.GOOS == "linux" {
		if cmd, err := startStdinPlayer(data); err == nil {
			return cmd, "", nil
		}
	}
	file, err := os.CreateTemp("", "yukumo-preview-*.wav")
	if err != nil {
		return nil, "", err
	}
	path := file.Name()
	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		_ = os.Remove(path)
		return nil, "", err
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(path)
		return nil, "", err
	}
	cmd, err := launchFilePlayer(path)
	if err != nil {
		_ = os.Remove(path)
		return nil, "", err
	}
	return cmd, path, nil
}

func startStdinPlayer(data []byte) (*exec.Cmd, error) {
	candidates := [][]string{
		{"paplay"},
		{"aplay", "-q"},
		{"ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet", "-"},
		{"mpv", "--no-video", "--really-quiet", "-"},
		{"play", "-q", "-"},
	}
	return startCandidate(data, candidates)
}

func launchFilePlayer(path string) (*exec.Cmd, error) {
	var candidates [][]string
	switch runtime.GOOS {
	case "darwin":
		candidates = [][]string{{"afplay", path}}
	case "windows":
		candidates = [][]string{
			{"powershell", "-NoProfile", "-Command",
				fmt.Sprintf("(New-Object Media.SoundPlayer '%s').PlaySync()", path)},
		}
	default:
		candidates = [][]string{
			{"paplay", path},
			{"aplay", "-q", path},
			{"ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet", path},
			{"mpv", "--no-video", "--really-quiet", path},
			{"play", "-q", path},
		}
	}
	return startCandidate(nil, candidates)
}

func startCandidate(data []byte, candidates [][]string) (*exec.Cmd, error) {
	var last error
	for _, args := range candidates {
		cmd := exec.Command(args[0], args[1:]...)
		if data != nil {
			cmd.Stdin = bytes.NewReader(data)
		}
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
		if err := cmd.Start(); err != nil {
			last = err
			continue
		}
		return cmd, nil
	}
	if last == nil {
		return nil, fmt.Errorf("preview: no audio player found")
	}
	return nil, fmt.Errorf("preview: %w", last)
}
