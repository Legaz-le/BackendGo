package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

const refreshTokenSize = 32

func GenerateRefreshToken() (rawToken string, tokenHash string, err error) {
	b := make([]byte, refreshTokenSize)
	if _, err = rand.Read(b); err != nil {
		return "", "", errors.New("Failed to generate refresh token: " + err.Error())
	}

	rawToken = base64.RawURLEncoding.EncodeToString(b)
	tokenHash = HashRefreshToken(rawToken)

	return rawToken, tokenHash, nil

}

func HashRefreshToken(rawToken string) string {
	sum := sha256.Sum256([]byte(rawToken))
	hashHex := hex.EncodeToString(sum[:])
	return hashHex
}
