package repository

import (
	"context"
	"testing"

	mockDatabase "go-unit-test/web-api/mocks/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
)

type RepositoryTest struct {
	suite.Suite
	mockMongoDatabase *mockDatabase.MongoDatabase
}

func (suite *RepositoryTest) SetupTest() {
	suite.mockMongoDatabase = &mockDatabase.MongoDatabase{}
}

func TestRepository(t *testing.T) {
	suite.Run(t, new(RepositoryTest))
}

func (suite *RepositoryTest) TestInsertOneUser() {
	suite.Run("Success", func() {
		insertedId := "asdasd"
		suite.mockMongoDatabase.On("InsertOne", mock.Anything, mock.Anything).Return(&insertedId, nil)

		repository := NewRepository(suite.mockMongoDatabase)
		_, err := repository.InsertOneUser(context.TODO(), bson.M{})
		assert.NoError(suite.T(), err)
	})
}

func (suite *RepositoryTest) TestFindOneUser() {
	suite.Run("Success", func() {
		suite.mockMongoDatabase.On("FindOne", mock.Anything, mock.Anything).Return(nil)

		repository := NewRepository(suite.mockMongoDatabase)
		_, err := repository.FindOneUser(context.TODO(), bson.M{})
		assert.NoError(suite.T(), err)
	})
}
