package scenario

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/arifsetiawan/json-db-benchmark/domain"
	"github.com/arifsetiawan/json-db-benchmark/pkg/random"
)

// NewDefinition ...
func NewDefinition() *domain.Definition {

	u, _ := uuid.NewV4()
	tenant := tenantID[RandomInt(0, 2)]
	referenceID := ""
	if tenant == "apple" {
		referenceID = referencesApple[RandomInt(0, 4)]
	} else {
		referenceID = referencesGoogle[RandomInt(0, 4)]
	}
	version := version[RandomInt(0, 3)]

	def := &domain.Definition{
		ID:        "definition:golang:" + u.String(),
		TenantID:  tenant,
		Entity:    "definition",
		Type:      "bpmn",
		CreatedAt: time.Now(),
		Reference: referenceID,
		Version:   version,
		Content:   random.GenerateAlphaNumeric(256),
		Elements:  make([]domain.Element, 0),
	}

	elementCount := RandomInt(5, 20)
	for i := 0; i < elementCount; i++ {
		def.Elements = append(def.Elements, domain.Element{
			ID:   random.GenerateNumeric(4),
			Name: random.GenerateAlphabet(12),
		})
	}

	return def
}

// NewInstance ...
func NewInstance(definition *domain.Definition) *domain.Instance {

	u, _ := uuid.NewV4()

	ins := &domain.Instance{
		ID:         "instance:golang:" + u.String(),
		TenantID:   definition.TenantID,
		Entity:     "instance",
		Type:       definition.Type,
		Definition: definition.ID,
		Reference:  definition.Reference,
		Version:    definition.Version,
		CreatedAt:  time.Now(),
		Activities: make([]domain.Activity, 0),
	}

	status := instanceStatus[RandomInt(0, 3)]
	switch status {
	case "completed":
		ins.Status.Completed = true
	case "failed":
		ins.Status.Failed = true
	case "running":
		ins.Status.Running = true
	}

	for _, v := range definition.Elements {
		ins.Activities = append(ins.Activities, domain.Activity{
			ID:        v.ID,
			Name:      v.Name,
			CreatedAt: time.Now(),
		})
	}

	return ins
}
