package auth

import (
	"encoding/json"
	"net/http"

	"github.com/goaperture/goaperture/v2/exception"
)

const (
	refreshPath = "/auth/refresh"
)

func (a *Auth[Payload]) BindHanders(server *http.ServeMux) {
	server.HandleFunc("/auth/login", a.onLogin)
	server.HandleFunc("/auth/logout", a.onLogout)
	server.HandleFunc(refreshPath, a.onRefresh)
}

type LoginInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginOutput struct {
	Token string `json:"accessToken"`
}

func (a *Auth[Payload]) onLogin(w http.ResponseWriter, r *http.Request) {
	defer exception.Catch(&w)

	var input LoginInput
	json.NewDecoder(r.Body).Decode(&input)

	id := a.Login(input.Login, input.Password)
	payload := a.GetPayload(id)

	a.createRefreshToken(&w, id)

	accessToken := a.getAccessToken(payload)
	json.NewEncoder(w).Encode(LoginOutput{accessToken})
}

func (a *Auth[Payload]) onLogout(w http.ResponseWriter, r *http.Request) {
	defer exception.Catch(&w)

}

func (a *Auth[Payload]) onRefresh(w http.ResponseWriter, r *http.Request) {
	defer exception.Catch(&w)

}
