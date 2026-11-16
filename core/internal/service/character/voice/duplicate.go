package voice

import (
	"context"
	"fmt"
	"strings"

	"yukumo-script/core/internal/domain"
	"yukumo-script/core/internal/service/nameutil"
)

func (s *Service) Duplicate(ctx context.Context, charName, srcName, newName string) (string, error) {
	srcName = strings.TrimSpace(srcName)
	if srcName == "" {
		return "", fmt.Errorf("%w: voice name is required", domain.ErrInvalid)
	}
	src, err := s.Get(charName, srcName)
	if err != nil {
		return "", err
	}
	dup := cloneVoice(src)
	dup.Name = nameutil.ResolveCopy(src.Name, newName)
	if err := s.Add(ctx, charName, dup); err != nil {
		return "", err
	}
	return dup.Name, nil
}

func cloneVoice(v domain.Voice) domain.Voice {
	out := v
	out.Params = make([]domain.SynthParam, len(v.Params))
	for i, p := range v.Params {
		out.Params[i] = p
		if len(p.PhontData) > 0 {
			out.Params[i].PhontData = append([]byte(nil), p.PhontData...)
		}
	}
	return out
}
