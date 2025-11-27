package utils

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

func NewID() (string, error) {
	var b [16]byte // 128-bit ID
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}

func Now() time.Time {
	return time.Now()
}
