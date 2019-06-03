package couchbase

import (
	"fmt"
	"os"

	"github.com/couchbase/gocb"

	"github.com/arifsetiawan/json-db-benchmark/domain"
)

// DefinitionStore is
type DefinitionStore struct {
	Bucket *gocb.Bucket
}

// NewDefinitionStore is
func NewDefinitionStore(bucket *gocb.Bucket) *DefinitionStore {
	return &DefinitionStore{
		Bucket: bucket,
	}
}

// Drop is
func (c *DefinitionStore) Drop() (err error) {

	bucketManager := c.Bucket.Manager(os.Getenv("COUCHBASE_ADMIN_USERNAME"), os.Getenv("COUCHBASE_ADMIN_PASSWORD"))
	err = bucketManager.Flush()
	if err != nil {
		return err
	}

	return nil
}

// Initialize is
func (c *DefinitionStore) Initialize() (err error) {

	bucketManager := c.Bucket.Manager(os.Getenv("COUCHBASE_ADMIN_USERNAME"), os.Getenv("COUCHBASE_ADMIN_PASSWORD"))
	bucketManager.CreatePrimaryIndex("engine_primary_index", true, false)

	bucketManager.CreateIndex("entity_tenant_id", []string{"entity", "tenant_id", "id"}, true, false)

	bucketManager.CreateIndex("entity_tenant_created_at", []string{"entity", "tenant_id", "-created_at"}, true, false)

	bucketManager.CreateIndex("entity_tenant_definition", []string{"entity", "tenant_id", "definition"}, true, false)

	return nil
}

// CreateDefinition is
func (c *DefinitionStore) CreateDefinition(data *domain.Definition) (err error) {

	if _, err := c.Bucket.Upsert(data.TenantID+":"+data.ID, data, 0); err != nil {
		return err
	}

	return nil
}

// GetDefinition is
func (c *DefinitionStore) GetDefinition(tenant string, id string) (data *domain.Definition, err error) {

	/*
		d := new(domain.Definition)
		if _, err := c.Bucket.Get(tenant+":"+id, &d); err != nil {
			return nil, err
		}

		return d, nil
	*/

	queryStr := fmt.Sprintf("SELECT t.* FROM engine t WHERE entity='definition' AND tenant_id='%s' AND id='%s'", tenant, id)

	query := gocb.NewN1qlQuery(queryStr)

	rows, err := c.Bucket.ExecuteN1qlQuery(query, nil)
	if err != nil {
		return nil, err
	}

	var item domain.Definition
	var items []domain.Definition
	for i := 0; rows.Next(&item); i++ {
		items = append(items, item)
		item = domain.Definition{}
	}
	_ = rows.Close()

	return &items[0], nil
}

// GetDefinitionWithStats is
func (c *DefinitionStore) GetDefinitionWithStats(tenant string, id string) (data *domain.Definition, err error) {

	query := gocb.NewN1qlQuery(`
	SELECT t.id, t.reference, t.version, running_count[0].count AS running_count, failed_count[0].count AS failed_count, completed_count[0].count AS completed_count
	FROM engine t 
	LET running_count = (
			SELECT count(*) count FROM engine e WHERE REGEXP_LIKE(META(e).id, '` + tenant + `:instance:.*') AND e.definition = '` + id + `' AND e.status.running = true
		),
		failed_count = (
			SELECT count(*) count FROM engine e WHERE REGEXP_LIKE(META(e).id, '` + tenant + `:instance:.*') AND e.definition = '` + id + `' AND e.status.failed = true
		),
		completed_count = (
			SELECT count(*) count FROM engine e WHERE REGEXP_LIKE(META(e).id, '` + tenant + `:instance:.*') AND e.definition = '` + id + `' AND e.status.completed = true
		)
		WHERE 
		REGEXP_LIKE(META(t).id, '` + tenant + `:definition:.*') AND 
		t.tenant_id IN ['` + tenant + `'] 
		AND t.id IN ['` + id + `']
	`)

	rows, err := c.Bucket.ExecuteN1qlQuery(query, nil)
	if err != nil {
		return nil, err
	}

	var item domain.Definition
	var items []domain.Definition
	for i := 0; rows.Next(&item); i++ {
		items = append(items, item)
		item = domain.Definition{}
	}
	_ = rows.Close()

	return &items[0], nil
}

// ListDefinitions is
func (c *DefinitionStore) ListDefinitions(tenant string, limit int64, offset int64) (data []domain.Definition, err error) {

	queryStr := fmt.Sprintf("SELECT t.* FROM engine t WHERE entity='definition' AND tenant_id='%s' ORDER BY created_at DESC LIMIT %d OFFSET %d", tenant, limit, offset)
	query := gocb.NewN1qlQuery(queryStr)

	rows, err := c.Bucket.ExecuteN1qlQuery(query, nil)
	if err != nil {
		return nil, err
	}

	var item domain.Definition
	var items []domain.Definition
	for i := 0; rows.Next(&item); i++ {
		items = append(items, item)
		item = domain.Definition{}
	}
	_ = rows.Close()

	return items, nil
}

// ListDefinitionsAll is
func (c *DefinitionStore) ListDefinitionsAll(limit int64, offset int64) (data []domain.Definition, err error) {

	queryStr := fmt.Sprintf("SELECT t.* FROM engine t WHERE entity='definition' ORDER BY created_at DESC LIMIT %d OFFSET %d", limit, offset)
	query := gocb.NewN1qlQuery(queryStr)

	rows, err := c.Bucket.ExecuteN1qlQuery(query, nil)
	if err != nil {
		return nil, err
	}

	var item domain.Definition
	var items []domain.Definition
	for i := 0; rows.Next(&item); i++ {
		items = append(items, item)
		item = domain.Definition{}
	}
	_ = rows.Close()

	return items, nil
}

// DeleteAll is
func (c *DefinitionStore) DeleteAll() (err error) {

	queryStr := fmt.Sprintf("DELETE FROM engine WHERE entity='definition'")
	query := gocb.NewN1qlQuery(queryStr)

	_, err = c.Bucket.ExecuteN1qlQuery(query, nil)
	if err != nil {
		return err
	}

	return nil
}
