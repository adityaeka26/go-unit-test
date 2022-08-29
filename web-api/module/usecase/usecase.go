package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"go-unit-test/web-api/database"
	"go-unit-test/web-api/helper"
	"go-unit-test/web-api/jwt"
	"go-unit-test/web-api/kafka"
	"go-unit-test/web-api/logger"
	"go-unit-test/web-api/module/model/domain"
	"go-unit-test/web-api/module/model/web"
	"go-unit-test/web-api/module/repository"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Usecase interface {
	Register(ctx context.Context, request web.RegisterRequest) error
	VerifyRegister(ctx context.Context, request web.VerifyRegisterRequest) (*web.VerifyRegisterResponse, error)
}
type UsecaseImpl struct {
	repository    repository.Repository
	jwtAuth       jwt.JWT
	redis         database.RedisClient
	kafkaProducer kafka.KafkaProducer
	logger        logger.Logger
}

func NewUsecase(repository repository.Repository, jwtAuth jwt.JWT, redis database.RedisClient, kafkaProducer kafka.KafkaProducer, logger logger.Logger) Usecase {
	return &UsecaseImpl{
		repository:    repository,
		jwtAuth:       jwtAuth,
		redis:         redis,
		kafkaProducer: kafkaProducer,
		logger:        logger,
	}
}

func (service UsecaseImpl) Register(ctx context.Context, request web.RegisterRequest) error {
	context := "service-Register"
	marshaledRequest, err := json.Marshal(request)
	if err != nil {
		service.logger.GetLogger().Error("Marshal request fail", zap.Error(err), zap.String("context", context))
		return helper.CustomError(http.StatusInternalServerError, "Internal server error")
	}

	// Preprocecss username & mobile number
	request.Username = strings.ToLower(request.Username)
	mobileNumberRegex := regexp.MustCompile(`^(\+628|628|08|8)`)
	if !mobileNumberRegex.MatchString(request.MobileNumber) {
		service.logger.GetLogger().Warn("Wrong format mobile number", zap.String("context", context), zap.ByteString("request", marshaledRequest))
		return helper.CustomError(http.StatusBadRequest, "Wrong format mobile number")
	}
	request.MobileNumber = mobileNumberRegex.ReplaceAllString(request.MobileNumber, "+628")

	// Check username is already registered or not yet
	user, err := service.repository.FindOneUser(ctx, bson.M{
		"username": request.Username,
	})
	if err != nil {
		service.logger.GetLogger().Error("Find one user fail", zap.Error(err), zap.String("context", context), zap.ByteString("request", marshaledRequest))
		return helper.CustomError(http.StatusInternalServerError, "Internal server error")
	}
	if user != nil {
		service.logger.GetLogger().Warn("Username already registered", zap.String("context", context), zap.ByteString("request", marshaledRequest))
		return helper.CustomError(http.StatusUnprocessableEntity, "Username already registered")
	}

	// Check username is already registered or not yet
	user, err = service.repository.FindOneUser(ctx, bson.M{
		"mobileNumber": request.MobileNumber,
	})
	if err != nil {
		service.logger.GetLogger().Error("Find one user 2 fail", zap.Error(err), zap.String("context", context), zap.ByteString("request", marshaledRequest))
		return helper.CustomError(http.StatusInternalServerError, "Internal server error")
	}
	if user != nil {
		service.logger.GetLogger().Warn("Mobile number already registered", zap.String("context", context), zap.ByteString("request", marshaledRequest))
		return helper.CustomError(http.StatusUnprocessableEntity, "Mobile number already registered")
	}

	// Generate hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		service.logger.GetLogger().Error("Generate hashed password fail", zap.Error(err), zap.String("context", context), zap.ByteString("request", marshaledRequest))
		return helper.CustomError(http.StatusInternalServerError, "Internal server error")
	}

	// Generate otp 6 digit
	otp, err := helper.GenerateOTP(6)
	if err != nil {
		service.logger.GetLogger().Error("Generate otp fail", zap.Error(err), zap.String("context", context), zap.ByteString("request", marshaledRequest))
		return helper.CustomError(http.StatusInternalServerError, "Internal server error")
	}

	// Generate user json for redis
	userRedis := domain.UserRedis{
		Username:     request.Username,
		Password:     string(hashedPassword),
		Name:         request.Name,
		Otp:          *otp,
		MobileNumber: request.MobileNumber,
	}
	userRedisJson, err := json.Marshal(userRedis)
	if err != nil {
		service.logger.GetLogger().Error("Marshal redis fail", zap.Error(err), zap.String("context", context), zap.ByteString("request", marshaledRequest))
		return helper.CustomError(http.StatusInternalServerError, "Internal server error")
	}

	// Insert user json to redis
	if err = service.redis.Set(
		ctx,
		fmt.Sprintf("REGISTER:%s", request.Username),
		string(userRedisJson),
		time.Minute*10,
	).Err(); err != nil {
		service.logger.GetLogger().Error("Set redis fail", zap.Error(err), zap.String("context", context), zap.ByteString("request", marshaledRequest))
		return helper.CustomError(http.StatusInternalServerError, "Internal server error")
	}

	// Generate user json for kafka
	userKafka := domain.RegisterOtpKafka{
		Name:         request.Name,
		Otp:          *otp,
		MobileNumber: request.MobileNumber,
	}
	userKafkaJson, err := json.Marshal(userKafka)
	if err != nil {
		service.logger.GetLogger().Error("Marshal kafka fail", zap.Error(err), zap.String("context", context), zap.ByteString("request", marshaledRequest))
		return helper.CustomError(http.StatusInternalServerError, "Internal server error")
	}

	// Produce user json to kafka
	err = service.kafkaProducer.SendMessage("REGISTER-OTP", string(userKafkaJson))
	if err != nil {
		service.logger.GetLogger().Error("Send message kafka fail", zap.Error(err), zap.String("context", context), zap.ByteString("request", marshaledRequest))
		return helper.CustomError(http.StatusInternalServerError, "Internal server error")
	}

	service.logger.GetLogger().Info(
		"Register success",
		zap.String("context", context),
		zap.ByteString("request", marshaledRequest),
	)

	return nil
}

func (service UsecaseImpl) VerifyRegister(ctx context.Context, request web.VerifyRegisterRequest) (*web.VerifyRegisterResponse, error) {
	context := "service-VerifyRegister"

	marshaledRequest, err := json.Marshal(request)
	if err != nil {
		service.logger.GetLogger().Error(
			"Marshal request fail",
			zap.String("context", context),
			zap.Error(err),
		)
		return nil, helper.CustomError(http.StatusInternalServerError, "Internal server error")
	}

	request.Username = strings.ToLower(request.Username)

	// Check username is already registered or not yet
	user, err := service.repository.FindOneUser(ctx, bson.M{
		"username": request.Username,
	})
	if err != nil {
		service.logger.GetLogger().Error(
			"Find one user fail",
			zap.String("context", context),
			zap.Error(err),
			zap.ByteString("request", marshaledRequest),
		)
		return nil, helper.CustomError(http.StatusInternalServerError, "Internal server error")
	}
	if user != nil {
		service.logger.GetLogger().Warn(
			"Username already registered",
			zap.String("context", context),
			zap.ByteString("request", marshaledRequest),
		)
		return nil, helper.CustomError(http.StatusUnprocessableEntity, "Username already registered")
	}

	// Get user from redis
	userRedisJson, err := service.redis.Get(ctx, fmt.Sprintf("REGISTER:%s", request.Username)).Result()
	if err != nil {
		if err == redis.Nil {
			service.logger.GetLogger().Warn(
				"OTP has expired",
				zap.String("context", context),
				zap.ByteString("request", marshaledRequest),
			)
			return nil, helper.CustomError(http.StatusUnauthorized, "OTP has expired")
		}
		service.logger.GetLogger().Error(
			"Get redis fail",
			zap.String("context", context),
			zap.Error(err),
			zap.ByteString("request", marshaledRequest),
		)
		return nil, helper.CustomError(http.StatusInternalServerError, "Internal server error")
	}

	// Parse user json to user object
	var userRedis domain.UserRedis
	err = json.Unmarshal([]byte(userRedisJson), &userRedis)
	if err != nil {
		service.logger.GetLogger().Error(
			"Unmarshal redis fail",
			zap.String("context", context),
			zap.Error(err),
			zap.ByteString("request", marshaledRequest),
		)
		return nil, helper.CustomError(http.StatusInternalServerError, "Internal server error")
	}

	// Validate otp
	if userRedis.Otp != request.Otp {
		service.logger.GetLogger().Warn(
			"Invalid OTP",
			zap.String("context", context),
			zap.ByteString("request", marshaledRequest),
		)
		return nil, helper.CustomError(http.StatusUnauthorized, "Invalid OTP")
	}

	// Insert user to mongodb
	insertUser := domain.InsertUser{
		Username:     userRedis.Username,
		Password:     userRedis.Password,
		Name:         userRedis.Name,
		MobileNumber: userRedis.MobileNumber,
	}
	insertedId, err := service.repository.InsertOneUser(ctx, insertUser)
	if err != nil {
		service.logger.GetLogger().Error(
			"Insert one user fail",
			zap.String("context", context),
			zap.Error(err),
			zap.ByteString("request", marshaledRequest),
		)
		return nil, helper.CustomError(http.StatusInternalServerError, "Internal server error")
	}

	// Delete user in redis
	err = service.redis.Del(ctx, fmt.Sprintf("REGISTER:%s", request.Username)).Err()
	if err != nil {
		service.logger.GetLogger().Error(
			"Delete redis fail",
			zap.String("context", context),
			zap.Error(err),
			zap.ByteString("request", marshaledRequest),
		)
		return nil, helper.CustomError(http.StatusInternalServerError, "Internal server error")
	}

	// Generate token
	token, err := service.jwtAuth.GenerateToken(jwt.Payload{
		Id: *insertedId,
	})
	if err != nil {
		service.logger.GetLogger().Error(
			"Generate token fail",
			zap.String("context", context),
			zap.Error(err),
			zap.ByteString("request", marshaledRequest),
		)
		return nil, helper.CustomError(http.StatusInternalServerError, "Internal server error")
	}

	response := web.VerifyRegisterResponse{
		Token: *token,
	}
	marshaledResponse, err := json.Marshal(response)
	if err != nil {
		service.logger.GetLogger().Error(
			"Marshal response fail",
			zap.String("context", context),
			zap.Error(err),
			zap.ByteString("request", marshaledRequest),
		)
		return nil, helper.CustomError(http.StatusInternalServerError, "Internal server error")
	}
	service.logger.GetLogger().Info(
		"Verify register success",
		zap.String("context", context),
		zap.ByteString("request", marshaledRequest),
		zap.ByteString("response", marshaledResponse),
	)

	return &response, nil
}
