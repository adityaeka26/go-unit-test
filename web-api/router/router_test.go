package router

import (
	"testing"

	"go-unit-test/web-api/middleware"
	mockHandler "go-unit-test/web-api/mocks/module/handler"

	"github.com/stretchr/testify/suite"
)

type RouterTest struct {
	suite.Suite
	mockHandler *mockHandler.Handler
	auth        middleware.Auth
}

func (suite *RouterTest) SetupTest() {
	suite.mockHandler = &mockHandler.Handler{}
	suite.auth = middleware.NewAuth()
}

func TestRouter(t *testing.T) {
	suite.Run(t, new(RouterTest))
}

func (suite *RouterTest) TestGetGinEngine() {
	suite.Run("Success", func() {
		router := NewRouter(suite.mockHandler, suite.auth)
		_ = router.GetGinEngine()
	})
}
