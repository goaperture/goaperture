package aperture

import (
	"net/http"
	"strings"
	"time"
)

type Payload interface {
	any
}

type Client[P Payload] struct {
	Payload *P
	Request
}

type Request struct {
	Request  *http.Request
	Responce *http.ResponseWriter
	secret   string
}

func NewClient[P Payload](r *http.Request, w *http.ResponseWriter, secret string) Client[P] {
	getPayload := func() *P {
		auth := r.Header.Get("Authorization")
		tokenString, _ := strings.CutPrefix(auth, "Bearer ")
		if tokenString == "" {
			return nil
		}

		payload, _ := DecodeAccessToken[P](tokenString, secret)

		return payload
	}

	return Client[P]{
		Payload: getPayload(),
		Request: Request{
			Request:  r,
			Responce: w,
			secret:   secret,
		},
	}
}

func (client Client[P]) NewJwt(payload P, refreshCookieKey string, refreshId string, path string, sequre bool) (string, error) {
	expires := time.Now().Add(7 * 24 * time.Hour)

	refresh, err := NewRefreshToken(refreshId, client.secret, expires)
	if err != nil {
		return "", err
	}

	http.SetCookie(*client.Responce, &http.Cookie{
		Name:     refreshCookieKey,
		Value:    refresh,
		Path:     path,
		Expires:  expires,
		HttpOnly: true,
		Secure:   sequre,
		SameSite: http.SameSiteStrictMode,
	})

	access, err := NewAccessToken(payload, client.secret)
	if err != nil {
		return "", err
	}

	return access, nil
}

func (client Client[P]) RefreshJwt(refreshCookieKey string, getPayload func(refreshId string) P) (string, error) {
	refreshToken, err := client.Request.Request.Cookie(refreshCookieKey)
	if err != nil {
		return "", err
	}

	clientId, err := DecodeRefreshToken(refreshToken.Value, client.secret)
	if err != nil {
		return "", err
	}

	payload := getPayload(clientId)

	return NewAccessToken(payload, client.secret)
}

func (client Client[P]) RemoveJwt(refreshCookieKey string, path string, sequre bool) {
	http.SetCookie(*client.Responce, &http.Cookie{
		Name:     refreshCookieKey,
		Value:    "",
		MaxAge:   -1,
		Path:     path,
		HttpOnly: true,
		Secure:   sequre,
		SameSite: http.SameSiteStrictMode,
	})
}
