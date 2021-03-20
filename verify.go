package main

import (
	"crypto/ed25519"
	"encoding/hex"
)

func verify(signature, hash, publicKey string) bool {
	decodedSignature, err := hex.DecodeString(signature)

	if err != nil {
		return false
	}

	decodedPublicKey, err := hex.DecodeString(publicKey)

	if err != nil {
		return false
	}

	return ed25519.Verify(decodedPublicKey, []byte(hash), decodedSignature)
}
