package synth

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// playWAV plays a WAV file with a platform-appropriate external player.
func playWAV(path string) error {
	var cmds [][]string
	switch runtime.GOOS {
	case "darwin":
		cmds = [][]string{{"afplay", path}}
	case "windows":
		cmds = [][]string{
			{"powershell", "-NoProfile", "-Command",
				fmt.Sprintf("(New-Object Media.SoundPlayer '%s').PlaySync()", path)},
		}
	default:
		cmds = [][]string{
			{"ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet", path},
			{"mpv", "--no-video", "--really-quiet", path},
			{"paplay", path},
			{"aplay", path},
			{"play", "-q", path},
		}
	}
	var last error
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			return nil
		} else {
			last = err
		}
	}
	if last == nil {
		return fmt.Errorf("no audio player found")
	}
	return fmt.Errorf("play wav: %w", last)
}
