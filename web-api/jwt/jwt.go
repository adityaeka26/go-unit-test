package jwt

import (
	"go-unit-test/web-api/config"

	"github.com/golang-jwt/jwt/v4"
)

type JWT interface {
	GenerateToken(payload Payload) (*string, error)
}

type JWTImpl struct {
	config config.Config
}

func NewJWT(config config.Config) JWT {
	return JWTImpl{
		config: config,
	}
}

type Payload struct {
	Id string `json:"id"`
}

func (jwtAuth JWTImpl) GenerateToken(payload Payload) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": payload.Id,
	})

	tokenString, err := token.SignedString([]byte(jwtAuth.config.GetEnv().JwtKey))
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}
