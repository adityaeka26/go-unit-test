package usecase

import (
	"context"
	"testing"

	"go-unit-test/basic/domain"
	"go-unit-test/basic/repository"
	mocksRepository "go-unit-test/basic/repository/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		_ = repository.NewRepository()
		mockRepository := &mocksRepository.Repository{}

		mockRepository.On("FindOneUser", mock.Anything, mock.Anything).Return(&domain.User{
			Name: "adit",
			Age:  25,
		}, nil)

		usecase := NewUsecase(mockRepository)
		actual, err := usecase.GetUser(context.TODO(), "asdasd")
		assert.NoError(t, err)
		assert.Equal(t, &domain.User{
			Name: "adit",
			Age:  25,
		}, actual)
	})

	t.Run("ErrorUserNotFound", func(t *testing.T) {
		_ = repository.NewRepository()
		mockRepository := &mocksRepository.Repository{}

		mockRepository.On("FindOneUser", mock.Anything, mock.Anything).Return(nil, nil)

		usecase := NewUsecase(mockRepository)
		actual, err := usecase.GetUser(context.TODO(), "asdasd")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}
