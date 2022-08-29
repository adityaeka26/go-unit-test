package usecase

import (
	"context"
	"testing"

	"go-unit-test/web-api/config"
	"go-unit-test/web-api/jwt"
	"go-unit-test/web-api/logger"
	mockDatabase "go-unit-test/web-api/mocks/database"
	mockJwt "go-unit-test/web-api/mocks/jwt"
	mockKafka "go-unit-test/web-api/mocks/kafka"
	mockRepository "go-unit-test/web-api/mocks/module/repository"
	"go-unit-test/web-api/module/model/domain"
	"go-unit-test/web-api/module/model/web"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UsecaseTest struct {
	suite.Suite
	config            config.Config
	logger            logger.Logger
	jwtAuth           jwt.JWT
	mockJwtAuth       *mockJwt.JWT
	mockRedisClient   *mockDatabase.RedisClient
	mockRepository    *mockRepository.Repository
	mockKafkaProducer *mockKafka.KafkaProducer
}

func (suite *UsecaseTest) SetupTest() {
	suite.config = config.NewConfig()
	suite.logger = logger.NewLog(suite.config.GetEnv().AppName, suite.config.GetEnv().AppEnvironment, config.NewConfig().GetEnv().LogPath)
	suite.jwtAuth = jwt.NewJWT(suite.config)
	suite.mockJwtAuth = &mockJwt.JWT{}
	suite.mockRedisClient = &mockDatabase.RedisClient{}
	suite.mockRepository = &mockRepository.Repository{}
	suite.mockKafkaProducer = &mockKafka.KafkaProducer{}
}

func TestUsecase(t *testing.T) {
	suite.Run(t, new(UsecaseTest))
}

func (suite *UsecaseTest) TestRegister() {
	suite.Run("Success", func() {
		suite.mockRepository.On("FindOneUser", mock.Anything, mock.Anything).Return(nil, nil).Once()
		suite.mockRepository.On("FindOneUser", mock.Anything, mock.Anything).Return(nil, nil).Once()
		suite.mockRedisClient.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(redis.NewStatusResult("sip", nil)).Once()
		suite.mockKafkaProducer.On("SendMessage", mock.Anything, mock.Anything).Return(nil).Once()

		usecase := NewUsecase(suite.mockRepository, suite.jwtAuth, suite.mockRedisClient, suite.mockKafkaProducer, suite.logger)
		err := usecase.Register(context.TODO(), web.RegisterRequest{
			MobileNumber: "08123123123",
		})
		assert.NoError(suite.T(), err)
	})

	suite.Run("ErrorWrongNumber", func() {
		usecase := NewUsecase(suite.mockRepository, suite.jwtAuth, suite.mockRedisClient, suite.mockKafkaProducer, suite.logger)
		err := usecase.Register(context.TODO(), web.RegisterRequest{
			MobileNumber: "123123123",
		})
		assert.Error(suite.T(), err)
	})

	suite.Run("ErrorUsernameAlreadyRegistered", func() {
		suite.mockRepository.On("FindOneUser", mock.Anything, mock.Anything).Return(&domain.User{}, nil).Once()

		usecase := NewUsecase(suite.mockRepository, suite.jwtAuth, suite.mockRedisClient, suite.mockKafkaProducer, suite.logger)
		err := usecase.Register(context.TODO(), web.RegisterRequest{
			MobileNumber: "08123123123",
		})
		assert.Error(suite.T(), err)
	})

	suite.Run("ErrorMobileNumberAlreadyRegistered", func() {
		suite.mockRepository.On("FindOneUser", mock.Anything, mock.Anything).Return(nil, nil).Once()
		suite.mockRepository.On("FindOneUser", mock.Anything, mock.Anything).Return(&domain.User{}, nil).Once()

		usecase := NewUsecase(suite.mockRepository, suite.jwtAuth, suite.mockRedisClient, suite.mockKafkaProducer, suite.logger)
		err := usecase.Register(context.TODO(), web.RegisterRequest{
			MobileNumber: "08123123123",
		})
		assert.Error(suite.T(), err)
	})
}

func (suite *UsecaseTest) TestVerifyRegister() {
	suite.Run("Success", func() {
		suite.mockRepository.On("FindOneUser", mock.Anything, mock.Anything).Return(nil, nil).Once()
		suite.mockRedisClient.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(redis.NewStringResult("{\"otp\":\"123123\"}", nil)).Once()
		insertedId := "123123123"
		suite.mockRepository.On("InsertOneUser", mock.Anything, mock.Anything).Return(&insertedId, nil).Once()
		suite.mockRedisClient.On("Del", mock.Anything, mock.Anything).Return(redis.NewIntResult(1, nil)).Once()
		token := "asdasd"
		suite.mockJwtAuth.On("GenerateToken", mock.Anything).Return(&token, nil)

		usecase := NewUsecase(suite.mockRepository, suite.mockJwtAuth, suite.mockRedisClient, suite.mockKafkaProducer, suite.logger)
		response, err := usecase.VerifyRegister(context.TODO(), web.VerifyRegisterRequest{
			Username: "08123123123",
			Otp:      "123123",
		})
		assert.Equal(suite.T(), &web.VerifyRegisterResponse{
			Token: token,
		}, response)
		assert.NoError(suite.T(), err)
	})

	suite.Run("ErrorUsernameAlreadyRegistered", func() {
		mockUser := domain.User{}
		suite.mockRepository.On("FindOneUser", mock.Anything, mock.Anything).Return(&mockUser, nil).Once()

		usecase := NewUsecase(suite.mockRepository, suite.jwtAuth, suite.mockRedisClient, suite.mockKafkaProducer, suite.logger)
		response, err := usecase.VerifyRegister(context.TODO(), web.VerifyRegisterRequest{
			Username: "08123123123",
			Otp:      "123123",
		})
		assert.Nil(suite.T(), response)
		assert.Error(suite.T(), err)
	})

	suite.Run("ErrorOtpExpired", func() {
		suite.mockRepository.On("FindOneUser", mock.Anything, mock.Anything).Return(nil, nil).Once()
		suite.mockRedisClient.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(redis.NewStringResult("", redis.Nil)).Once()

		usecase := NewUsecase(suite.mockRepository, suite.mockJwtAuth, suite.mockRedisClient, suite.mockKafkaProducer, suite.logger)
		response, err := usecase.VerifyRegister(context.TODO(), web.VerifyRegisterRequest{
			Username: "08123123123",
			Otp:      "123123",
		})
		assert.Nil(suite.T(), response)
		assert.Error(suite.T(), err)
	})

	suite.Run("ErrorInvalidOtp", func() {
		suite.mockRepository.On("FindOneUser", mock.Anything, mock.Anything).Return(nil, nil).Once()
		suite.mockRedisClient.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(redis.NewStringResult("{\"otp\":\"123123\"}", nil)).Once()

		usecase := NewUsecase(suite.mockRepository, suite.mockJwtAuth, suite.mockRedisClient, suite.mockKafkaProducer, suite.logger)
		response, err := usecase.VerifyRegister(context.TODO(), web.VerifyRegisterRequest{
			Username: "08123123123",
			Otp:      "321321",
		})
		assert.Nil(suite.T(), response)
		assert.Error(suite.T(), err)
	})
}
