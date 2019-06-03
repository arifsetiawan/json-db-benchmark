package postgres

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"

	"github.com/arifsetiawan/json-db-benchmark/domain"
)

// InstanceStore is
type InstanceStore struct {
	db *sqlx.DB
}

// NewInstanceStore is
func NewInstanceStore(db *sqlx.DB) *InstanceStore {
	return &InstanceStore{
		db: db,
	}
}

// Drop is
func (c *InstanceStore) Drop() (err error) {
	stmt := `
		DROP TABLE IF EXISTS engine_instance CASCADE;
	`

	_, err = c.db.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}

// Initialize is
func (c *InstanceStore) Initialize() (err error) {
	stmt := `
		CREATE TABLE  IF NOT EXISTS engine_instance (
			instance_id varchar(64) PRIMARY KEY,
			tenant_id varchar(16) NOT NULL,
			definition_id varchar(64) NOT NULL,
			data jsonb NOT NULL,
			created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX json_instance_status ON engine_instance USING gin (data);
		
		CREATE INDEX instance_def_id ON engine_instance (definition_id);
	`

	_, err = c.db.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}

// CreateInstance is
func (c *InstanceStore) CreateInstance(data *domain.Instance) (err error) {

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	stmt := `
		INSERT INTO engine_instance (instance_id, tenant_id, definition_id, data)
		VALUES ($1, $2, $3, $4);
	`

	_, err = c.db.Exec(stmt, data.ID, data.TenantID, data.Definition, string(b))
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
	_, err = c.db.Exec("DELETE from engine_instance")
	if err != nil {
		return err
	}

	return nil
}
