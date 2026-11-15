package synth

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"

	"yukumo-script/internal/domain"
)

// TaskHash fingerprints text, effective speed, synth params, and optional character-file bytes.
func TaskHash(text string, speed int, p domain.SynthParam, characterFile []byte) string {
	h := sha256.New()
	writeString(h, text)
	_ = binary.Write(h, binary.LittleEndian, int32(speed))
	_ = binary.Write(h, binary.LittleEndian, int32(p.Version))
	_ = binary.Write(h, binary.LittleEndian, int32(p.Preset))
	_ = binary.Write(h, binary.LittleEndian, int32(p.Speed))
	writeBytes(h, p.PhontData)
	_ = binary.Write(h, binary.LittleEndian, int32(p.Base))
	_ = binary.Write(h, binary.LittleEndian, int32(p.Volume))
	_ = binary.Write(h, binary.LittleEndian, int32(p.Pitch))
	_ = binary.Write(h, binary.LittleEndian, int32(p.Accent))
	_ = binary.Write(h, binary.LittleEndian, int32(p.Lmd))
	_ = binary.Write(h, binary.LittleEndian, int32(p.Fsc))
	writeBytes(h, characterFile)
	return hex.EncodeToString(h.Sum(nil))
}

type hashWriter interface {
	Write([]byte) (int, error)
}

func writeString(h hashWriter, s string) {
	_ = binary.Write(h, binary.LittleEndian, uint32(len(s)))
	_, _ = h.Write([]byte(s))
}

func writeBytes(h hashWriter, b []byte) {
	_ = binary.Write(h, binary.LittleEndian, uint32(len(b)))
	if len(b) > 0 {
		_, _ = h.Write(b)
	}
}
