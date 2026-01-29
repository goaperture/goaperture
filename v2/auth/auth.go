package auth

import (
	"net/http"
	"time"

	"github.com/goaperture/goaperture/v2/auth/auth_paths"
)

func (a *Auth[Payload]) GetSecret() XSecret {
	return XSecret{rsa: &a.RSA, strSecret: a.Secret}
}

func (a *Auth[Payload]) getAccessToken(client Payload) string {
	var token = getJwt(client, a.LiveTime.AccessKey, a.GetSecret())
	return token
}

func (a *Auth[Payload]) createRefreshToken(w *http.ResponseWriter, id ID) {
	life := a.LiveTime.RefreshKey

	if life == 0 {
		life = 24 * 60 * 7
	}

	var token = getJwt(PrivatePayload{id}, life, a.GetSecret())
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
