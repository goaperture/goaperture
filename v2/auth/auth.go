package auth

import "github.com/goaperture/goaperture/v2/auth_rsa"

type ID string
type Permission string
type Permissions []Permission

type Secret struct {
	Key struct {
		Public  string
		Private string
	}
	SecretKey string
}

type LiveTime struct {
	AccessKey  int
	RefreshKey int
}

type RSA struct {
	Public  *auth_rsa.PublicPemKey
	Private *auth_rsa.PrivatePemKey
}

type Auth[Payload any] struct {
	Sequre     bool
	LiveTime   LiveTime
	Login      func(login, password string) ID
	GetPayload func(id ID) Payload
	Secret     Secret
	RSA        RSA
}
