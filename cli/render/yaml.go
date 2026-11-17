package render

import (
	"io"

	"gopkg.in/yaml.v3"
)

// YAML writes v as indented YAML.
func YAML(w io.Writer, v any) error {
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	defer func() { _ = enc.Close() }()
	return enc.Encode(v)
}
