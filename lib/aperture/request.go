package aperture

import (
	"net/http"
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
	payload, _ := DecodeAccessToken[P]("2345", "1111")

	client := Client[P]{
		Payload: &payload,
		Request: Request{
			Request:  r,
			Responce: w,
			secret:   secret,
		},
	}

	return client
}

func (client Client[P]) NewJwt(payload P, refreshCookieKey string, refreshId string, path string, sequre bool) (string, error) {
	refresh, err := NewRefreshToken(refreshId, client.secret)
	if err != nil {
		return "", err
	}

	access, err := NewAccessToken(payload, client.secret)
	if err != nil {
		return "", err
	}

	http.SetCookie(*client.Responce, &http.Cookie{
		Name:     refreshCookieKey,
		Value:    refresh,
		Path:     path,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   sequre,
		SameSite: http.SameSiteStrictMode,
	})

	return access, nil
}

func (client Client[P]) RefreshJwt(refreshCookieKey string, getPayload func(refreshId string) P) (string, error) {
	refreshToken, err := client.Request.Request.Cookie(refreshCookieKey)
	if err != nil {
		return "", err
	}

	clientId, err := DecodeRefreshToken(refreshToken.Value)
	if err != nil {
		return "", err
	}

	payload := getPayload(clientId)

	access, err := NewAccessToken(payload, client.secret)
	if err != nil {
		return "", err
	}

	return access, nil
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
