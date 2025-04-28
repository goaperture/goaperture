package aperture

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Client any `json:"payload"`
	jwt.RegisteredClaims
}

func NewAccessToken(client Payload, secret string) (string, error) {
	claims := CustomClaims{
		Client: client,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(secret)
}

func NewRefreshToken(clientId string, secret string) (string, error) {
	return "123", nil
}

func DecodeAccessToken[P Payload](tokenString string, secret string) (P, error) {
	var payload P

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return payload, err
	}

	if !token.Valid {
		return payload, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return payload, errors.New("invalid claims type")
	}

	data, err := json.Marshal(claims)
	if err != nil {
		return payload, err
	}

	err = json.Unmarshal(data, &payload)
	if err != nil {
		return payload, err
	}

	return payload, nil
}

func DecodeRefreshToken(token string) (clientId string, err error) {

	return
}
