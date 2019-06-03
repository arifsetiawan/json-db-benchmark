package couchbase

import (
	"fmt"

	"github.com/couchbase/gocb"

	"github.com/arifsetiawan/json-db-benchmark/domain"
)

// InstanceStore is
type InstanceStore struct {
	Bucket *gocb.Bucket
}

// NewInstanceStore is
func NewInstanceStore(bucket *gocb.Bucket) *InstanceStore {
	return &InstanceStore{
		Bucket: bucket,
	}
}

// Drop is
func (c *InstanceStore) Drop() (err error) {
	// do nothing. dropped in definition
	return nil
}

// Initialize is
func (c *InstanceStore) Initialize() (err error) {
	// do nothing. initialized in definition
	return nil
}

// CreateInstance is
func (c *InstanceStore) CreateInstance(data *domain.Instance) (err error) {

	if _, err := c.Bucket.Upsert(data.TenantID+":"+data.ID, data, 0); err != nil {
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

	queryStr := fmt.Sprintf("DELETE FROM engine WHERE entity='instance'")
	query := gocb.NewN1qlQuery(queryStr)

	_, err = c.Bucket.ExecuteN1qlQuery(query, nil)
	if err != nil {
		return err
	}

	return nil
}
