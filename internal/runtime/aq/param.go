package aq

import "encoding/binary"

// Param mirrors third-party/aq_wasm.h AqSynthParam (40-byte LE block).
type Param struct {
	Version int32
	Preset  int32
	Speed   int32
	Phont   []byte // V2 embedded .phont; empty means none
	Base    int32
	Volume  int32
	Pitch   int32
	Accent  int32
	Lmd     int32
	Fsc     int32
}

const paramSize = 40

func encodeParam(p Param, phontOff uint32) []byte {
	b := make([]byte, paramSize)
	binary.LittleEndian.PutUint32(b[0:4], uint32(p.Version))
	binary.LittleEndian.PutUint32(b[4:8], uint32(p.Preset))
	binary.LittleEndian.PutUint32(b[8:12], uint32(p.Speed))
	binary.LittleEndian.PutUint32(b[12:16], phontOff)
	binary.LittleEndian.PutUint32(b[16:20], uint32(p.Base))
	binary.LittleEndian.PutUint32(b[20:24], uint32(p.Volume))
	binary.LittleEndian.PutUint32(b[24:28], uint32(p.Pitch))
	binary.LittleEndian.PutUint32(b[28:32], uint32(p.Accent))
	binary.LittleEndian.PutUint32(b[32:36], uint32(p.Lmd))
	binary.LittleEndian.PutUint32(b[36:40], uint32(p.Fsc))
	return b
}
