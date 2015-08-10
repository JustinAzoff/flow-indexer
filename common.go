package main

import (
	"encoding/binary"
)

func PutUVarint(v uint64) []byte {
	b := make([]byte, 10)
	binary.PutUvarint(b, uint64(v))
	return b
}
