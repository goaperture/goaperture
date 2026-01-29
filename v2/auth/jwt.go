package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/goaperture/goaperture/v2/auth/auth_paths"
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

func (a *Auth[Payload]) getAccessToken(client Payload) string {
	var token = getJwt(client, a.LiveTime.AccessKey)
	return token
}

func (a *Auth[Payload]) createRefreshToken(w *http.ResponseWriter, id ID) {
	life := a.LiveTime.RefreshKey

	if life == 0 {
		life = 24 * 60 * 7
	}

	var token = getJwt(PrivatePayload{id}, life)
	expires := time.Now().Add(time.Minute * time.Duration(life))

	http.SetCookie(*w, &http.Cookie{
		Name:     refreshCookieKey,
		Value:    token,
		Path:     auth_paths.REFRESH,
		Expires:  expires,
		HttpOnly: true,
		Secure:   a.Sequre,
		SameSite: http.SameSiteStrictMode,
	})

}

func (a *Auth[Payload]) removeRefreshToken(w *http.ResponseWriter) {
	http.SetCookie(*w, &http.Cookie{
		Name:     refreshCookieKey,
		Value:    "",
		MaxAge:   -1,
		Path:     auth_paths.REFRESH,
		HttpOnly: true,
		Secure:   a.Sequre,
		SameSite: http.SameSiteStrictMode,
	})
}

func getRefreshToken(r *http.Request) string {
	cookie, err := r.Cookie(refreshCookieKey)
	if err != nil {
		exception.Fall("Не удалось получить refresh token", "invalid Refresh Token", 401)
	}

	return cookie.Value
}

func getJwt[T any](payload T, life int) string {
	secret := ""

	if life == 0 {
		life = 5
	}

	claims := CustomClaims[T]{
		Payload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(life))),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	result, err := token.SignedString([]byte(secret))
	if err != nil {
		exception.Fall("Не удалось создать token", "fall Signed Token", 401)
	}

	return result
}

func GetPayloadFromJwt[P any](tokenString string) *P {
	var secret = ""

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims[P]{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
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
