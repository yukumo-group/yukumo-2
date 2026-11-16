package prompt

import "github.com/charmbracelet/huh"

var accessible bool

func SetAccessible(v bool) { accessible = v }

func Accessible() bool { return accessible }

func run(field huh.Field) error {
	form := huh.NewForm(huh.NewGroup(field)).WithShowHelp(false)
	if accessible {
		form = form.WithAccessible(true)
	}
	return form.Run()
}
