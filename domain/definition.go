package domain

import "time"

// DefinitionStorage ...
type DefinitionStorage interface {
	Drop() (err error)
	Initialize() (err error)

	CreateDefinition(data *Definition) (err error)

	GetDefinition(tenant string, id string) (data *Definition, err error)
	GetDefinitionWithStats(tenant string, id string) (data *Definition, err error)

	ListDefinitions(tenant string, limit int64, offset int64) (data []Definition, err error)

	ListDefinitionsAll(limit int64, offset int64) (data []Definition, err error)

	DeleteAll() (err error)
}

// Definition ...
type Definition struct {
	ID             string    `json:"id" bson:"id" db:"id"`
	TenantID       string    `json:"tenant_id" bson:"tenant_id" db:"tenant_id"`
	Entity         string    `json:"entity" bson:"entity" db:"entity"`
	Type           string    `json:"type" bson:"type" db:"type"`
	Reference      string    `json:"reference" bson:"reference" db:"reference"`
	Version        string    `json:"version" bson:"version" db:"version"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at" db:"created_at"`
	CompletedCount int       `json:"completed_count" bson:"completed_count" db:"completed_count"`
	FailedCount    int       `json:"failed_count" bson:"failed_count" db:"failed_count"`
	RunningCount   int       `json:"running_count" bson:"running_count" db:"running_count"`
	Content        string    `json:"content" bson:"content" db:"content"`
	Elements       []Element `json:"elements" bson:"elements" db:"elements"`
}

// Element ...
type Element struct {
	ID   string `json:"id" bson:"id" bson:"id"`
	Name string `json:"name" bson:"name" bson:"name"`
}
