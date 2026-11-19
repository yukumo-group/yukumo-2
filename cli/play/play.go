package play

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"yukumo-script/core/api"
)

// Preview plays the synthesized WAV with a platform-appropriate external player.
func Preview(_ context.Context, res api.SynthResult) error {
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
	defer func() {
		if cleanup != "" {
			_ = os.Remove(cleanup)
		}
	}()
	return playWAV(path)
}

func playWAV(path string) error {
	cmd, err := launchPlayer(path)
	if err != nil {
		return err
	}
	return cmd.Wait()
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
