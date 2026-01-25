package auth

import (
	"net/http"
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

func (a *Auth[Payload]) getAccessToken(client Payload) string {
	var token = getJwt(client, a.LiveTime.AccessKey)
	return token
}

func (a *Auth[Payload]) createRefreshToken(w *http.ResponseWriter, id ID) {
	life := a.LiveTime.RefreshKey

	if life == 0 {
		life = 24 * 60
	}

	var token = getJwt(PrivatePayload{id}, a.LiveTime.RefreshKey)
	expires := time.Now().Add(time.Minute * time.Duration(a.LiveTime.RefreshKey))

	http.SetCookie(*w, &http.Cookie{
		Name:     refreshCookieKey,
		Value:    token,
		Path:     refreshPath,
		Expires:  expires,
		HttpOnly: true,
		Secure:   a.Sequre,
		SameSite: http.SameSiteStrictMode,
	})

}

func getJwt[T any](payload T, life int) string {
	secret := ""

	if life == 0 {
		life = 10
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
		exception.Fall("Не удалось создать Токен", "4511")
	}

	return result
}
