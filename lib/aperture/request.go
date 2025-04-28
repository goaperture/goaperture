package aperture

import (
	"net/http"
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

func (client Client[P]) NewAccessToken(payload P) (string, error) {
	return NewAccessToken(payload, client.secret)
}
