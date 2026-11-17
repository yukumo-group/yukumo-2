package domain

type Character struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description,omitempty"`
	Voices      []Voice `yaml:"voices,omitempty"`
}
