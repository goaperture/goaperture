package auth

type ID string

type LiveTime struct {
	AccessKey  int
	RefreshKey int
}

type Auth[Payload any] struct {
	Sequre     bool
	LiveTime   LiveTime
	Login      func(login, password string) ID
	GetPayload func(id ID) Payload
	Secret     string
	RSA        RSA
}
