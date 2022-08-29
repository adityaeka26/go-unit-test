package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-unit-test/web-api/config"
	"go-unit-test/web-api/logger"
	mockUsecase "go-unit-test/web-api/mocks/module/usecase"
	"go-unit-test/web-api/module/model/web"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type HandlerTest struct {
	suite.Suite
	config      config.Config
	mockUsecase *mockUsecase.Usecase
	logger      logger.Logger
}

func (suite *HandlerTest) SetupTest() {
	suite.config = config.NewConfig()
	suite.mockUsecase = &mockUsecase.Usecase{}
	suite.logger = logger.NewLog(suite.config.GetEnv().AppName, config.NewConfig().GetEnv().AppEnvironment, config.NewConfig().GetEnv().LogPath)
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(HandlerTest))
}

func (suite *HandlerTest) TestRegister() {
	suite.Run("Success", func() {
		suite.mockUsecase.On("Register", mock.Anything, mock.Anything).Return(nil)

		handler := NewHandler(suite.mockUsecase, suite.logger)
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		body := web.RegisterRequest{
			Username:     "asd",
			Password:     "asd",
			Name:         "asd",
			MobileNumber: "08123123123",
		}
		bodyJson, _ := json.Marshal(body)
		c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyJson))
		c.Request.Header.Add("Content-Type", "application/json")
		handler.Register(c)
	})
}

func (suite *HandlerTest) TestVerifyRegister() {
	suite.Run("Success", func() {
		suite.mockUsecase.On("VerifyRegister", mock.Anything, mock.Anything).Return(nil, nil)

		handler := NewHandler(suite.mockUsecase, suite.logger)
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		body := web.VerifyRegisterRequest{
			Username: "asdasd",
			Otp:      "123123",
		}
		bodyJson, _ := json.Marshal(body)
		c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyJson))
		c.Request.Header.Add("Content-Type", "application/json")
		handler.VerifyRegister(c)
	})
}
