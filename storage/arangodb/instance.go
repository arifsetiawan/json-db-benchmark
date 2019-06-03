package arangodb

import (
	"context"

	arangodb "github.com/arangodb/go-driver"

	"github.com/arifsetiawan/json-db-benchmark/domain"
)

// InstanceStore is
type InstanceStore struct {
	collection     arangodb.Collection
	database       arangodb.Database
	collectionName string
}

// NewInstanceStore is
func NewInstanceStore(database arangodb.Database, collectionName string) *InstanceStore {
	return &InstanceStore{
		database:       database,
		collectionName: collectionName,
		collection:     nil,
	}
}

// Drop is
func (c *InstanceStore) Drop() (err error) {

	exists, err := c.database.CollectionExists(context.Background(), c.collectionName)
	if err != nil {
		return err
	}

	if exists {
		c.collection, err = c.database.Collection(context.Background(), c.collectionName)
		if err != nil {
			return err
		}

		err = c.collection.Remove(context.Background())
		if err != nil {
			return err
		}
	}

	return nil
}

// Initialize is
func (c *InstanceStore) Initialize() (err error) {

	c.collection, err = c.database.CreateCollection(context.Background(), c.collectionName, nil)
	if err != nil {
		return err
	}

	_, _, err = c.collection.EnsureSkipListIndex(context.Background(), []string{"definition"}, nil)
	if err != nil {
		return err
	}

	_, _, err = c.collection.EnsureSkipListIndex(context.Background(), []string{"status.completed"}, nil)
	if err != nil {
		return err
	}

	_, _, err = c.collection.EnsureSkipListIndex(context.Background(), []string{"status.running"}, nil)
	if err != nil {
		return err
	}

	_, _, err = c.collection.EnsureSkipListIndex(context.Background(), []string{"status.failed"}, nil)
	if err != nil {
		return err
	}

	return nil
}

// CreateInstance is
func (c *InstanceStore) CreateInstance(data *domain.Instance) (err error) {
	_, err = c.collection.CreateDocument(context.Background(), data)
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

// ListInstancesByReference is
func (c *InstanceStore) ListInstancesByReference(tenant string, reference string, version string) (data []domain.Instance, err error) {
	return nil, nil
}

// DeleteAll is
func (c *InstanceStore) DeleteAll() (err error) {
	ctx := context.Background()
	query := "FOR d IN instances REMOVE d IN instances"
	cursor, err := c.database.Query(ctx, query, nil)
	if err != nil {
		return err
	}
	defer cursor.Close()

	return nil
}
