package auth

import (
	"crypto/rsa"
	"fmt"

	"github.com/goaperture/goaperture/v2/auth_rsa"
	"github.com/goaperture/goaperture/v2/exception"
	"github.com/golang-jwt/jwt/v5"
)

type RSA struct {
	Public  *auth_rsa.PublicPemKey
	Private *auth_rsa.PrivatePemKey
}

type xsecret struct {
	rsa       *RSA
	strSecret string
}

func createTokenWithPrivateSign[T any](claims CustomClaims[T], privateKey *auth_rsa.PrivatePemKey) string {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	key := (*rsa.PrivateKey)(*privateKey)
	result, err := token.SignedString(key)
	if err != nil {
		exception.Fall(fmt.Sprintf("%s", err), "fall Signed Token", 401)
	}

	return result
}

func createTokenWithSecret[T any](claims CustomClaims[T], secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	result, err := token.SignedString([]byte(secret))
	if err != nil {
		exception.Fall("Не удалось подписать token", "fall Signed Token", 401)
	}

	return result
}
