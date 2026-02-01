package auth

import (
	"encoding/json"
	"net/http"

	"github.com/goaperture/goaperture/v2/api/auth/auth_paths"
	"github.com/goaperture/goaperture/v2/exception"
	"github.com/goaperture/goaperture/v2/responce"
)

func (a *Auth[Payload]) BindHanders(server *http.ServeMux) {
	server.HandleFunc(auth_paths.LOGIN, a.onLogin)
	server.HandleFunc(auth_paths.LOGOUT, a.onLogout)
	server.HandleFunc(auth_paths.REFRESH, a.onRefresh)
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
	// var input = params.GetInput[LoginInput](r)

	id := a.Login(input.Login, input.Password)
	payload := a.GetPayload(id)

	a.createRefreshToken(&w, id)

	accessToken := a.getAccessToken(payload)
	json.NewEncoder(w).Encode(LoginOutput{accessToken})
}

func (a *Auth[Payload]) onLogout(w http.ResponseWriter, r *http.Request) {
	defer exception.Catch(&w)

	a.removeRefreshToken(&w)

	json.NewEncoder(w).Encode(responce.Success(true))
}

func (a *Auth[Payload]) onRefresh(w http.ResponseWriter, r *http.Request) {
	defer exception.Catch(&w)

	var refreshToken = getRefreshToken(r)
	privatePayload := GetPayloadFromJwt[PrivatePayload](refreshToken, a.GetSecret())

	payload := a.GetPayload(privatePayload.Id)
	accessToken := a.getAccessToken(payload)
	json.NewEncoder(w).Encode(LoginOutput{accessToken})

}
