package voice

import (
	"fmt"
	"strconv"

	"yukumo-script/cli/action"
	"yukumo-script/cli/render"
	"yukumo-script/cli/streams"
	"yukumo-script/core/api"
)

func runEdit(s streams.Streams, app *api.App, charName string, args []string) error {
	if len(args) < 1 || args[0] == "" {
		return fmt.Errorf("voice name is required")
	}
	if len(args) < 2 {
		return fmt.Errorf("field and value required")
	}
	return runOneShotVoiceEdit(s, app, charName, args[0], args[1:])
}

func runOneShotVoiceEdit(s streams.Streams, app *api.App, charName, voiceName string, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("field and value required")
	}
	field, raw := args[0], args[1]
	switch field {
	case "name":
		_, result, err := action.ApplyVoiceName(s.Ctx, app, charName, voiceName, raw)
		if err != nil {
			return err
		}
		return render.Line(s.Out, "name: %s\n", result.Value)
	case "speed":
		_, result, err := action.ApplyVoiceSpeed(s.Ctx, app, charName, voiceName, raw, true)
		if err != nil {
			return err
		}
		if result.Warning != "" {
			_ = render.Warning(s.Err, result.Warning)
		}
		return render.Line(s.Out, "speed: %s\n", result.Value)
	case "param", "version":
		return fmt.Errorf("use: edit <voice> <version> <field> <value>")
	default:
		if len(args) >= 3 && (field == "1" || field == "2" || field == "10") {
			ver, _ := strconv.Atoi(field)
			if _, result, err := action.ApplyVoiceParam(s.Ctx, app, charName, voiceName, ver, args[1], args[2], true); err != nil {
				return err
			} else if result.Warning != "" {
				_ = render.Warning(s.Err, result.Warning)
			}
			return render.Line(s.Out, "v%d.%s updated\n", ver, args[1])
		}
		return fmt.Errorf("unknown field %q (want name, speed, or <version> <field> <value>)", field)
	}
}
