package repository

import (
	"context"
	"go-unit-test/basic/domain"
)

type Repository interface {
	FindOneUser(ctx context.Context, filter interface{}) (*domain.User, error)
}

type RepositoryImpl struct{}

func NewRepository() Repository {
	return &RepositoryImpl{}
}

func (repository *RepositoryImpl) FindOneUser(ctx context.Context, filter interface{}) (*domain.User, error) {
	return nil, nil
}
