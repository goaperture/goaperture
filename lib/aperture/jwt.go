package aperture

import (
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Client any `json:"payload"`
	jwt.RegisteredClaims
}

type CustomRefreshClaims struct {
	ClientId string `json:"client-id"`
	jwt.RegisteredClaims
}

func getSecretKey(secret string) []byte {
	return []byte(secret)
}

func NewAccessToken(client Payload, secret string) (string, error) {
	claims := CustomClaims{
		Client: client,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecretKey(secret))
}

func NewRefreshToken(clientId string, secret string, lifeTime time.Time) (string, error) {
	claims := CustomRefreshClaims{
		ClientId: clientId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(lifeTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecretKey(secret))
}

func DecodeAccessToken[P Payload](tokenString string, secret string) (*P, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return getSecretKey(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if climps, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		data, err := json.Marshal(climps.Client)

		if err != nil {
			return nil, err
		}

		var payload P
		err = json.Unmarshal(data, &payload)
		if err == nil {
			return &payload, err
		}
	}

	return nil, err
}

func DecodeRefreshToken(tokenString string, secret string) (clientId string, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomRefreshClaims{}, func(t *jwt.Token) (interface{}, error) {
		return getSecretKey(secret), nil
	})

	if climps, ok := token.Claims.(*CustomRefreshClaims); ok && token.Valid {
		clientId = climps.ClientId
	}

	return
}
