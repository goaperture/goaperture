package auth

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/goaperture/goaperture/v2/exception"
	"github.com/golang-jwt/jwt/v5"
)

const (
	refreshCookieKey = "RefreshToken"
)

type PrivatePayload struct {
	Id ID `json:"id"`
}

type CustomClaims[T any] struct {
	Payload T `json:"payload"`
	jwt.RegisteredClaims
}

func getRefreshToken(r *http.Request) string {
	cookie, err := r.Cookie(refreshCookieKey)
	if err != nil {
		exception.Fall("Не удалось получить refresh token", "invalid Refresh Token", 401)
	}

	return cookie.Value
}

func getJwt[T any](payload T, life int, secret XSecret) string {
	if life == 0 {
		life = 5
	}

	claims := CustomClaims[T]{
		Payload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(life))),
		},
	}

	if secret.rsa.Private != nil {
		return createTokenWithPrivateSign(claims, secret.rsa.Private)
	} else {
		return createTokenWithSecret(claims, secret.strSecret)
	}

}

func GetPayloadFromJwt[P any](tokenString string, secret XSecret) *P {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims[P]{}, func(t *jwt.Token) (interface{}, error) {
		if secret.rsa.Public != nil {
			key := (*rsa.PublicKey)(*secret.rsa.Public)
			return key, nil
		}

		return []byte(secret.strSecret), nil
	})

	if err != nil {
		exception.Fall(fmt.Sprintf("%s", err), "invalid token", 401)
	}

	if v, ok := token.Claims.(*CustomClaims[P]); ok && token.Valid {
		return &v.Payload
	}

	return nil
}

func ParseAccessToken(r *http.Request) (string, bool) {
	auth := r.Header.Get("Authorization")
	return strings.CutPrefix(auth, "Bearer ")
}
