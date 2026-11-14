package character_test

import (
	"errors"
	"testing"

	"yukumo-script/internal/domain"
	svccharacter "yukumo-script/internal/service/character"
)

func TestFileName(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr error
	}{
		{name: "Alice", want: "Alice.yaml"},
		{name: "  Bob  ", want: "Bob.yaml"},
		{name: `a/b\c:d*e?f"g<h>i|j`, want: "a_b_c_d_e_f_g_h_i_j.yaml"},
		{name: "...", wantErr: domain.ErrInvalid},
		{name: "???", want: "___.yaml"},
		{name: "", wantErr: domain.ErrInvalid},
		{name: "   ", wantErr: domain.ErrInvalid},
	}
	for _, tt := range tests {
		got, err := svccharacter.FileName(tt.name)
		if tt.wantErr != nil {
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("FileName(%q) err=%v want %v", tt.name, err, tt.wantErr)
			}
			continue
		}
		if err != nil {
			t.Fatalf("FileName(%q): %v", tt.name, err)
		}
		if got != tt.want {
			t.Fatalf("FileName(%q)=%q want %q", tt.name, got, tt.want)
		}
	}
}
