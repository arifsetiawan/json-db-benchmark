package postgres

import (
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/arifsetiawan/json-db-benchmark/domain"
)

// DefinitionData ...
type DefinitionData struct {
	DefinitionID string     `db:"definition_id"`
	TenantID     string     `db:"tenant_id"`
	Data         string     `db:"data"`
	CreatedAt    *time.Time `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
}

// DefinitionStore is
type DefinitionStore struct {
	db *sqlx.DB
}

// NewDefinitionStore is
func NewDefinitionStore(db *sqlx.DB) *DefinitionStore {
	return &DefinitionStore{
		db: db,
	}
}

// Drop is
func (c *DefinitionStore) Drop() (err error) {
	stmt := `
		DROP TABLE IF EXISTS engine_definition CASCADE;
	`

	_, err = c.db.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}

// Initialize is
func (c *DefinitionStore) Initialize() (err error) {
	stmt := `
		CREATE TABLE IF NOT EXISTS engine_definition (
			definition_id varchar(64) PRIMARY KEY,
			tenant_id varchar(16) NOT NULL,
			data jsonb NOT NULL,
			created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX definition_tenant_id ON engine_definition (tenant_id, definition_id);
		
		CREATE INDEX definition_id ON engine_definition (definition_id);
		
		CREATE INDEX definition_tenant_created_at ON engine_definition (tenant_id, created_at DESC);
		
		CREATE INDEX json_definition_tenant_id ON engine_definition ((data->>'tenant_id'),(data->>'id'));
		
		CREATE INDEX json_definition_tenant_created_at ON engine_definition ((data->>'tenant_id'),(data->>'created_at') DESC);
	`

	_, err = c.db.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}

// CreateDefinition is
func (c *DefinitionStore) CreateDefinition(data *domain.Definition) (err error) {

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	stmt := `
		INSERT INTO engine_definition (definition_id, tenant_id, data)
		VALUES ($1, $2, $3);
	`

	_, err = c.db.Exec(stmt, data.ID, data.TenantID, string(b))
	if err != nil {
		return err
	}

	return nil
}

// GetDefinition is
func (c *DefinitionStore) GetDefinition(tenant string, id string) (data *domain.Definition, err error) {

	d := DefinitionData{}
	//stmt := "SELECT * FROM engine_definition WHERE tenant_id=$1 AND definition_id=$2"
	stmt := "SELECT * FROM engine_definition WHERE data->>'tenant_id'=$1 AND data->>'id'=$2"
	err = c.db.Get(&d, stmt, tenant, id)
	if err != nil {
		return nil, err
	}

	def := domain.Definition{}
	err = json.Unmarshal([]byte(d.Data), &def)
	if err != nil {
		return nil, err
	}

	return &def, nil
}

// GetDefinitionWithStats is
func (c *DefinitionStore) GetDefinitionWithStats(tenant string, id string) (data *domain.Definition, err error) {

	d := domain.Definition{}
	stmt := `SELECT 
		d.definition_id as id, 
		d.data->>'reference' as reference, 
		d.data->>'version' as version, 
		SUM(CASE WHEN i.data @> '{"status":{"running": true}}' THEN 1 ELSE 0 END) AS running_count,
		SUM(CASE WHEN i.data @> '{"status":{"completed": true}}' THEN 1 ELSE 0 END) AS completed_count,
		SUM(CASE WHEN i.data @> '{"status":{"failed": true}}' THEN 1 ELSE 0 END) AS failed_count
	FROM engine_definition d 
	FULL OUTER JOIN engine_instance i ON i.definition_id = d.definition_id
	WHERE
		d.data->>'tenant_id'=$1 AND
		d.data->>'id'=$2
	GROUP BY d.definition_id`

	err = c.db.Get(&d, stmt, tenant, id)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ListDefinitions is
func (c *DefinitionStore) ListDefinitions(tenant string, limit int64, offset int64) (data []domain.Definition, err error) {

	d := []DefinitionData{}
	err = c.db.Select(&d, "SELECT * FROM engine_definition WHERE data->>'tenant_id'=$1 ORDER BY data->>'created_at' DESC OFFSET $2 LIMIT $3", tenant, offset, limit)
	if err != nil {
		return nil, err
	}

	for _, v := range d {
		def := domain.Definition{}
		err := json.Unmarshal([]byte(v.Data), &def)
		if err != nil {
			return nil, err
		}
		data = append(data, def)
	}

	return data, nil
}

// ListDefinitionsAll is
func (c *DefinitionStore) ListDefinitionsAll(limit int64, offset int64) (data []domain.Definition, err error) {

	d := []DefinitionData{}
	err = c.db.Select(&d, "SELECT * FROM engine_definition ORDER BY data->>'created_at' DESC OFFSET $1 LIMIT $2", offset, limit)
	if err != nil {
		return nil, err
	}

	for _, v := range d {
		def := domain.Definition{}
		err := json.Unmarshal([]byte(v.Data), &def)
		if err != nil {
			return nil, err
		}
		data = append(data, def)
	}

	return data, nil
}

// DeleteAll is
func (c *DefinitionStore) DeleteAll() (err error) {
	_, err = c.db.Exec("DELETE from engine_definition")
	if err != nil {
		return err
	}

	return nil
}
