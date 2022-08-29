package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Auth interface {
	VerifyJWTToken() gin.HandlerFunc
}

type AuthImpl struct{}

func NewAuth() Auth {
	return &AuthImpl{}
}

func (auth AuthImpl) VerifyJWTToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("test")
	}
}
