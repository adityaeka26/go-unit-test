package usecase

import (
	"context"
	"go-unit-test/basic/domain"
	"go-unit-test/basic/repository"
	mocksRepository "go-unit-test/basic/repository/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UsecaseTestSuite struct {
	suite.Suite
	repository     repository.Repository
	mockRepository *mocksRepository.Repository
}

func (suite *UsecaseTestSuite) SetupTest() {
	suite.repository = repository.NewRepository()
	suite.mockRepository = &mocksRepository.Repository{}
}

func TestUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UsecaseTestSuite))
}

func (suite *UsecaseTestSuite) TestGetUserSuite() {
	suite.Run("Success", func() {
		suite.mockRepository.On("FindOneUser", mock.Anything, mock.Anything).Return(&domain.User{
			Name: "adit",
			Age:  25,
		}, nil).Once()

		usecase := NewUsecase(suite.mockRepository)
		actual, err := usecase.GetUser(context.TODO(), "asdasd")
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), &domain.User{
			Name: "adit",
			Age:  25,
		}, actual)
	})

	suite.Run("ErrorUserNotFound", func() {
		suite.mockRepository.On("FindOneUser", mock.Anything, mock.Anything).Return(nil, nil).Once()

		usecase := NewUsecase(suite.mockRepository)
		actual, err := usecase.GetUser(context.TODO(), "asdasd")
		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), actual)
	})
}
