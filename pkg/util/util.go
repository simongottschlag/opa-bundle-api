package util

import (
	"crypto/sha256"
	"fmt"
)

var (
	NullString = ""
)

func BytesToHash(data []byte) (string, error) {
	hasher := sha256.New()
	_, err := hasher.Write(data)
	if err != nil {
		return NullString, err
	}

	hash := fmt.Sprintf("%x", hasher.Sum(nil))
	return hash, nil
}
