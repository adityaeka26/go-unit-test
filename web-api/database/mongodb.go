package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDatabase interface {
	InsertOne(ctx context.Context, payload InsertOne) (*string, error)
	FindOne(ctx context.Context, payload FindOne) error
}

type MongoDatabaseImpl struct {
	mongoClient  *mongo.Client
	databaseName string
}

func NewMongoDB(uri string, dbName string) MongoDatabase {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	return &MongoDatabaseImpl{
		mongoClient:  client,
		databaseName: dbName,
	}
}

type InsertOne struct {
	CollectionName string
	Document       interface{}
}

func (mongoDatabase MongoDatabaseImpl) InsertOne(ctx context.Context, payload InsertOne) (*string, error) {
	collection := mongoDatabase.mongoClient.Database(mongoDatabase.databaseName).Collection(payload.CollectionName)
	insertDoc, err := collection.InsertOne(ctx, payload.Document)
	if err != nil {
		return nil, err
	}

	insertedId := insertDoc.InsertedID.(primitive.ObjectID).Hex()
	return &insertedId, nil
}

type FindOne struct {
	CollectionName string
	Filter         interface{}
	Result         interface{}
}

func (mongoDatabase MongoDatabaseImpl) FindOne(ctx context.Context, payload FindOne) error {
	collection := mongoDatabase.mongoClient.Database(mongoDatabase.databaseName).Collection(payload.CollectionName)
	result := collection.FindOne(ctx, payload.Filter)

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil
		}
		return result.Err()
	}

	if err := result.Decode(payload.Result); err != nil {
		return err
	}

	return nil
}
