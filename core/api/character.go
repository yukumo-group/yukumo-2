package api

func (a *App) ListCharacters() ([]Character, error) {
	list, err := a.chars.List()
	if err != nil {
		return nil, err
	}
	out := make([]Character, len(list))
	for i, c := range list {
		out[i] = toCharacter(c)
	}
	return out, nil
}

func (a *App) GetCharacter(name string) (Character, error) {
	c, err := a.chars.Get(name)
	if err != nil {
		return Character{}, err
	}
	return toCharacter(c), nil
}

func (a *App) AddCharacter(c Character) error {
	return a.chars.Add(fromCharacter(c))
}

func (a *App) UpdateCharacter(oldName string, c Character) error {
	return a.chars.Update(oldName, fromCharacter(c))
}

func (a *App) DeleteCharacter(name string) error {
	return a.chars.Delete(name)
}
