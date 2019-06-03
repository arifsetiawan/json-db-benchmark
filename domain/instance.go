package domain

import "time"

// InstanceStorage ...
type InstanceStorage interface {
	Drop() (err error)
	Initialize() (err error)

	CreateInstance(data *Instance) (err error)

	GetInstance(tenant string, id string) (data *Instance, err error)

	ListInstances(tenant string) (data []Instance, err error)

	DeleteAll() (err error)
}

// Instance ...
type Instance struct {
	ID         string     `json:"id" bson:"id" db:"id"`
	TenantID   string     `json:"tenant_id" bson:"tenant_id" db:"tenant_id"`
	Entity     string     `json:"entity" bson:"entity" db:"entity"`
	Type       string     `json:"type" bson:"type" db:"type"`
	Definition string     `json:"definition" bson:"definition" db:"definition"`
	Reference  string     `json:"reference" bson:"reference" db:"reference"`
	Version    string     `json:"version" bson:"version" db:"version"`
	CreatedAt  time.Time  `json:"created_at" bson:"created_at" db:"created_at"`
	Status     Status     `json:"status" bson:"status" db:"status"`
	Activities []Activity `json:"activities" bson:"activities" db:"activities"`
}

// Status ...
type Status struct {
	Completed bool `json:"completed" bson:"completed" db:"completed"`
	Failed    bool `json:"failed" bson:"failed" db:"failed"`
	Running   bool `json:"running" bson:"running" db:"running"`
}

// Activity ...
type Activity struct {
	ID        string    `json:"id" bson:"id" db:"id"`
	Name      string    `json:"name" bson:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" bson:"created_at" db:"created_at"`
}
