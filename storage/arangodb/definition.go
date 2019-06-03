package arangodb

import (
	"context"

	arangodb "github.com/arangodb/go-driver"

	"github.com/arifsetiawan/json-db-benchmark/domain"
)

// DefinitionStore is
type DefinitionStore struct {
	collection     arangodb.Collection
	database       arangodb.Database
	collectionName string
}

// NewDefinitionStore is
func NewDefinitionStore(database arangodb.Database, collectionName string) *DefinitionStore {
	return &DefinitionStore{
		database:       database,
		collectionName: collectionName,
		collection:     nil,
	}
}

// Drop is
func (c *DefinitionStore) Drop() (err error) {

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
func (c *DefinitionStore) Initialize() (err error) {

	c.collection, err = c.database.CreateCollection(context.Background(), c.collectionName, nil)
	if err != nil {
		return err
	}

	_, _, err = c.collection.EnsureSkipListIndex(context.Background(), []string{"tenant_id", "id"}, nil)
	if err != nil {
		return err
	}

	_, _, err = c.collection.EnsureSkipListIndex(context.Background(), []string{"tenant_id", "created_at"}, nil)
	if err != nil {
		return err
	}

	return nil
}

// CreateDefinition is
func (c *DefinitionStore) CreateDefinition(data *domain.Definition) (err error) {

	_, err = c.collection.CreateDocument(context.Background(), data)
	if err != nil {
		return err
	}

	return nil
}

// GetDefinition is
func (c *DefinitionStore) GetDefinition(tenant string, id string) (data *domain.Definition, err error) {

	query := "FOR d IN definitions FILTER d.tenant_id == @tenant && d.id == @id LIMIT 1 RETURN d"
	bindVars := map[string]interface{}{
		"tenant": tenant,
		"id":     id,
	}
	cursor, err := c.database.Query(context.Background(), query, bindVars)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var definition domain.Definition
	for {
		_, err := cursor.ReadDocument(context.Background(), &definition)
		if arangodb.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}
	}

	return &definition, nil
}

// GetDefinitionWithStats is
func (c *DefinitionStore) GetDefinitionWithStats(tenant string, id string) (data *domain.Definition, err error) {

	query := `
		FOR d IN definitions
			FILTER d.tenant_id == @tenant && d.id == @id
			LET instances_running = (
				FOR f IN instances
					FILTER f.definition == d.id && f.status.running == true
				RETURN f
			)
			LET instances_completed = (
				FOR f IN instances
					FILTER f.definition == d.id && f.status.completed == true
				RETURN f
			)
			LET instances_failed = (
				FOR f IN instances
					FILTER f.definition == d.id && f.status.failed == true
				RETURN f
			)
			
		RETURN {
			id: d.id, 
			reference: d.reference, 
			version: d.version, 
			running_count: LENGTH(instances_running),
			completed_count: LENGTH(instances_completed),
			failed_count: LENGTH(instances_failed)
		}
	`
	bindVars := map[string]interface{}{
		"tenant": tenant,
		"id":     id,
	}
	cursor, err := c.database.Query(context.Background(), query, bindVars)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var definition domain.Definition
	for {
		_, err := cursor.ReadDocument(context.Background(), &definition)
		if arangodb.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}
	}

	return &definition, nil
}

// ListDefinitions is
func (c *DefinitionStore) ListDefinitions(tenant string, limit int64, offset int64) (data []domain.Definition, err error) {

	query := "FOR d IN definitions FILTER d.tenant_id == @tenant SORT d.created_at DESC LIMIT @offset, @limit RETURN d"
	bindVars := map[string]interface{}{
		"tenant": tenant,
		"limit":  limit,
		"offset": offset,
	}
	cursor, err := c.database.Query(context.Background(), query, bindVars)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var items []domain.Definition
	for {
		var item domain.Definition
		_, err := cursor.ReadDocument(context.Background(), &item)
		if arangodb.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

// ListDefinitionsAll is
func (c *DefinitionStore) ListDefinitionsAll(limit int64, offset int64) (data []domain.Definition, err error) {

	query := "FOR d IN definitions SORT d.created_at DESC LIMIT @offset, @limit RETURN d"
	bindVars := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}
	cursor, err := c.database.Query(context.Background(), query, bindVars)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var items []domain.Definition
	for {
		var item domain.Definition
		_, err := cursor.ReadDocument(context.Background(), &item)
		if arangodb.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

// DeleteAll is
func (c *DefinitionStore) DeleteAll() (err error) {
	ctx := context.Background()
	query := "FOR d IN definitions REMOVE d IN definitions"
	cursor, err := c.database.Query(ctx, query, nil)
	if err != nil {
		return err
	}
	defer cursor.Close()

	return nil
}
