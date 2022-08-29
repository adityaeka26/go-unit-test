package main

import (
	"go-unit-test/web-api/config"
	"go-unit-test/web-api/database"
	"go-unit-test/web-api/jwt"
	"go-unit-test/web-api/kafka"
	"go-unit-test/web-api/logger"
	"go-unit-test/web-api/middleware"
	"go-unit-test/web-api/module/handler"
	"go-unit-test/web-api/module/repository"
	"go-unit-test/web-api/module/usecase"
	"go-unit-test/web-api/router"
)

func main() {
	config := config.NewConfig()

	jwtAuth := jwt.NewJWT(config)
	authMiddleware := middleware.NewAuth()
	logger := logger.NewLog(config.GetEnv().AppName, config.GetEnv().AppEnvironment, config.GetEnv().LogPath)

	kafkaProducer := kafka.NewKafkaProducer(config.GetEnv().KafkaUrl, logger)
	redis := database.NewRedis(config.GetEnv().RedisUrl, config.GetEnv().RedisPassword)
	mongoDatabase := database.NewMongoDB(config.GetEnv().MongodbUrl, config.GetEnv().MongodbDatabaseName)

	repository := repository.NewRepository(mongoDatabase)
	usecase := usecase.NewUsecase(repository, jwtAuth, redis.GetClient(), kafkaProducer, logger)
	handler := handler.NewHandler(usecase, logger)
	router := router.NewRouter(handler, authMiddleware)

	router.GetGinEngine().Run(":8080")
}
