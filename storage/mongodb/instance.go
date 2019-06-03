package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"github.com/arifsetiawan/json-db-benchmark/domain"
)

// InstanceStore is
type InstanceStore struct {
	collection     *mongo.Collection
	database       *mongo.Database
	collectionName string
}

// NewInstanceStore is
func NewInstanceStore(database *mongo.Database, collectionName string) *InstanceStore {
	return &InstanceStore{
		collection:     database.Collection(collectionName),
		database:       database,
		collectionName: collectionName,
	}
}

// Drop is
func (c *InstanceStore) Drop() (err error) {
	err = c.database.RunCommand(
		context.Background(),
		bsonx.Doc{{"drop", bsonx.String(c.collectionName)}},
	).Err()
	if err != nil {
		return err
	}

	return nil
}

// Initialize is
func (c *InstanceStore) Initialize() (err error) {

	// create collection
	err = c.database.RunCommand(
		context.Background(),
		bsonx.Doc{{"create", bsonx.String(c.collectionName)}},
	).Err()
	if err != nil {
		return err
	}

	c.collection = c.database.Collection(c.collectionName)

	// index
	indexView := c.collection.Indexes()

	// create index
	defIndex := mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "definition", Value: bsonx.Int32(1)},
		},
	}
	statusRunnning := mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "status.running", Value: bsonx.Int32(1)},
		},
	}
	statusFailed := mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "status.failed", Value: bsonx.Int32(1)},
		},
	}
	statusCompleted := mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "status.completed", Value: bsonx.Int32(1)},
		},
	}

	_, err = indexView.CreateMany(context.Background(),
		[]mongo.IndexModel{defIndex, statusRunnning, statusFailed, statusCompleted})
	if err != nil {
		return err
	}

	return nil
}

// CreateInstance is
func (c *InstanceStore) CreateInstance(data *domain.Instance) (err error) {

	_, err = c.collection.InsertOne(context.Background(), data)
	if err != nil {
		return err
	}

	return nil
}

// GetInstance ...
func (c *InstanceStore) GetInstance(tenant string, id string) (data *domain.Instance, err error) {
	return nil, nil
}

// ListInstances is
func (c *InstanceStore) ListInstances(tenant string) (data []domain.Instance, err error) {
	return nil, nil
}

// DeleteAll is
func (c *InstanceStore) DeleteAll() (err error) {
	ctx := context.Background()
	_, err = c.collection.DeleteMany(
		ctx,
		bson.D{},
	)
	if err != nil {
		return err
	}

	return nil
}
