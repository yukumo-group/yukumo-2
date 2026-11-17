package store

import (
	"errors"

	"gopkg.in/yaml.v3"

	"yukumo-script/core/port"
)

func LoadYAML[T any](r port.Reader, name string) (T, error) {
	var zero T
	data, err := r.Read(name)
	if err != nil {
		return zero, err
	}
	var v T
	if err := yaml.Unmarshal(data, &v); err != nil {
		return zero, err
	}
	return v, nil
}

func SaveYAML[T any](w port.Writer, name string, v T) error {
	data, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	return w.Write(name, data)
}

func LoadYAMLOrDefault[T any](r port.Reader, name string, def T) (T, error) {
	v, err := LoadYAML[T](r, name)
	if errors.Is(err, port.ErrNotFound) {
		return def, nil
	}
	return v, err
}
