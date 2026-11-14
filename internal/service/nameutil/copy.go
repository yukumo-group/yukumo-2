package nameutil

import "strings"

const CopySuffix = " - copy"

// ResolveCopy returns newName when set, otherwise name + " - copy".
func ResolveCopy(name, newName string) string {
	name = strings.TrimSpace(name)
	newName = strings.TrimSpace(newName)
	if newName != "" {
		return newName
	}
	return name + CopySuffix
}
