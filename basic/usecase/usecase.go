package usecase

import (
	"context"
	"errors"

	"go-unit-test/basic/domain"
	"go-unit-test/basic/repository"

	"go.mongodb.org/mongo-driver/bson"
)

type Usecase interface {
	GetUser(ctx context.Context, id string) (*domain.User, error)
}

type UsecaseImpl struct {
	repository repository.Repository
}

func NewUsecase(repository repository.Repository) Usecase {
	return &UsecaseImpl{
		repository: repository,
	}
}

func (usecase *UsecaseImpl) GetUser(ctx context.Context, id string) (*domain.User, error) {
	user, err := usecase.repository.FindOneUser(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
