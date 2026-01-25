package auth

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
	PublicKey  string
	PrivateKey string

	publicToken  string
	privateToken string
}

type Auth[Payload any] struct {
	Sequre     bool
	LiveTime   LiveTime
	Login      func(login, password string) ID
	GetPayload func(id ID) Payload
	Secret     Secret
	RSA        RSA
}
