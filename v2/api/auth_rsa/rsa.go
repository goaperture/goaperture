package auth_rsa

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type PrivatePemKey *rsa.PrivateKey
type PublicPemKey *rsa.PublicKey

func LoadPrivateKey(path string) *PrivatePemKey {
	keyData, err := os.ReadFile(path)
	if err != nil {
		panic("Не удалось прочитать приватный ключ")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		panic("Не удалось распарсить приватный ключ")
	}

	keyLink := PrivatePemKey(privateKey)

	return &keyLink
}

func LoadPublicKey(path string) *PublicPemKey {
	keyData, err := os.ReadFile(path)
	if err != nil {
		panic("Не удалось прочитать публичный ключ")
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		panic("Не удалось распарсить публичный ключ")
	}

	keyLink := PublicPemKey(publicKey)

	return &keyLink
}
