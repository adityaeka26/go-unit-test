package repository

import (
	"context"

	"go-unit-test/web-api/database"
	"go-unit-test/web-api/module/model/domain"
)

type Repository interface {
	InsertOneUser(ctx context.Context, document interface{}) (*string, error)
	FindOneUser(ctx context.Context, filter interface{}) (*domain.User, error)
}
type RepositoryImpl struct {
	mongo database.MongoDatabase
}

func NewRepository(mongo database.MongoDatabase) Repository {
	return &RepositoryImpl{
		mongo: mongo,
	}
}

func (repository *RepositoryImpl) InsertOneUser(ctx context.Context, document interface{}) (*string, error) {
	insertedId, err := repository.mongo.InsertOne(ctx, database.InsertOne{
		CollectionName: "users",
		Document:       document,
	})
	if err != nil {
		return nil, err
	}

	return insertedId, nil
}

func (repository *RepositoryImpl) FindOneUser(ctx context.Context, filter interface{}) (*domain.User, error) {
	var result *domain.User
	err := repository.mongo.FindOne(ctx, database.FindOne{
		CollectionName: "users",
		Filter:         filter,
		Result:         &result,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
