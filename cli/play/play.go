package play

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync"

	"yukumo-script/core/api"
)

var (
	playMu  sync.Mutex
	playCmd *exec.Cmd
	playTmp string
)

// Preview plays the synthesized WAV. If async is true, playback starts in the
// background so interactive UIs can continue.
func Preview(_ context.Context, res api.SynthResult, async bool) error {
	path := res.Path
	cleanup := ""
	if path == "" {
		f, err := os.CreateTemp("", "yukumo-synth-*.wav")
		if err != nil {
			return err
		}
		cleanup = f.Name()
		if _, err := f.Write(res.WAV); err != nil {
			_ = f.Close()
			_ = os.Remove(cleanup)
			return err
		}
		if err := f.Close(); err != nil {
			_ = os.Remove(cleanup)
			return err
		}
		path = cleanup
	}
	if async {
		playWAVAsync(path, cleanup)
		return nil
	}
	defer func() {
		if cleanup != "" {
			_ = os.Remove(cleanup)
		}
	}()
	return playWAV(path)
}

// playWAV plays a WAV file with a platform-appropriate external player (blocking).
func playWAV(path string) error {
	cmd, err := launchPlayer(path)
	if err != nil {
		return err
	}
	return cmd.Wait()
}

// playWAVAsync starts playback in the background and does not block the UI.
// If cleanup is non-empty, that path is removed after playback ends (or is stopped).
// A new call stops any previous background playback.
func playWAVAsync(path, cleanup string) {
	go func() {
		stopPlayback()

		cmd, err := launchPlayer(path)
		if err != nil {
			if cleanup != "" {
				_ = os.Remove(cleanup)
			}
			return
		}

		playMu.Lock()
		playCmd = cmd
		playTmp = cleanup
		playMu.Unlock()

		_ = cmd.Wait()

		playMu.Lock()
		if playCmd == cmd {
			playCmd = nil
			if playTmp != "" {
				_ = os.Remove(playTmp)
				playTmp = ""
			}
		}
		playMu.Unlock()
	}()
}

func stopPlayback() {
	playMu.Lock()
	cmd := playCmd
	tmp := playTmp
	playCmd = nil
	playTmp = ""
	playMu.Unlock()

	if cmd != nil && cmd.Process != nil {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}
	if tmp != "" {
		_ = os.Remove(tmp)
	}
}

func launchPlayer(path string) (*exec.Cmd, error) {
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
			{"ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet", path},
			{"mpv", "--no-video", "--really-quiet", path},
			{"paplay", path},
			{"aplay", path},
			{"play", "-q", path},
		}
	}
	var last error
	for _, args := range candidates {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
		if err := cmd.Start(); err != nil {
			last = err
			continue
		}
		return cmd, nil
	}
	if last == nil {
		return nil, fmt.Errorf("no audio player found")
	}
	return nil, fmt.Errorf("play wav: %w", last)
}
