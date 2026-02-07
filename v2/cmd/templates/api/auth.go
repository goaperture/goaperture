package api

import (
	"github.com/goaperture/goaperture/v2/api/auth"
	"github.com/goaperture/goaperture/v2/exception"
)

type Payload struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Avatar           string `json:"avatar"`
	auth.Permissions `json:"permissions"`
}

var Auth = auth.Auth[Payload]{
	Login:      login,
	GetPayload: getPayload,
	// RSA: auth.RSA{
	// 	Public:  auth_rsa.LoadPublicKey("./rsa/public.pem"),
	// 	Private: auth_rsa.LoadPrivateKey("./rsa/private.pem"),
	// },
	// Secret: "SECRET-KEY-OPTION",
}

func login(login, password string) auth.ID {
	if login == "admin" && password == "admin" {
		return "111"
	}

	exception.Fall("not valid login or email", "Not Valid Email Or Password", 401)
	return ""
}

func getPayload(id auth.ID) Payload {
	return Payload{
		Id:          string(id),
		Name:        "Admin",
		Permissions: auth.Permissions{"hello_hello_world"},
	}
}
