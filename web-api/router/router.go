package router

import (
	"go-unit-test/web-api/middleware"
	"go-unit-test/web-api/module/handler"

	"github.com/gin-gonic/gin"
	"go.elastic.co/apm/module/apmgin"
)

type Router interface {
	GetGinEngine() *gin.Engine
}

type RouterImpl struct {
	ginEngine *gin.Engine
	auth      middleware.Auth
}

func NewRouter(handler handler.Handler, auth middleware.Auth) Router {
	router := gin.New()
	router.Use(apmgin.Middleware(router))

	userRouter := router.Group("/user")

	userRouter.Use(apmgin.Middleware(router))

	userRouter.GET("/", handler.Index)
	userRouter.POST("/v1/register", handler.Register)
	userRouter.POST("/v1/register/verify", handler.VerifyRegister)

	return &RouterImpl{
		ginEngine: router,
	}
}

func (router RouterImpl) GetGinEngine() *gin.Engine {
	return router.ginEngine
}
