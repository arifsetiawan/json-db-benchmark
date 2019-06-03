package scenario

import (
	"github.com/rs/zerolog/log"

	"github.com/arifsetiawan/json-db-benchmark/domain"
)

// Scenario ...
type Scenario struct {
	StorageType       string
	DefinitionStorage domain.DefinitionStorage
	InstanceStorage   domain.InstanceStorage
}

// Drop ...
func (c *Scenario) Drop() (err error) {

	err = c.DefinitionStorage.Drop()
	if err != nil {
		log.Error().Str("type", c.StorageType).Err(err).Msg("drop definitions")
	}

	err = c.InstanceStorage.Drop()
	if err != nil {
		log.Error().Str("type", c.StorageType).Err(err).Msg("drop instances")
	}

	return nil
}

// Initialize ...
func (c *Scenario) Initialize() (err error) {

	err = c.DefinitionStorage.Initialize()
	if err != nil {
		log.Error().Str("type", c.StorageType).Err(err).Msg("initialize definitions")
	}

	err = c.InstanceStorage.Initialize()
	if err != nil {
		log.Error().Str("type", c.StorageType).Err(err).Msg("initialize definitions")
	}

	return nil
}

// CreateSingleDefinition ...
func (c *Scenario) CreateSingleDefinition() (err error) {

	definition := NewDefinition()
	err = c.DefinitionStorage.CreateDefinition(definition)
	if err != nil {
		log.Error().Str("type", c.StorageType).Err(err).Msg("insert document")
	}

	return nil
}

// CreateDefinitionAndInstances ...
func (c *Scenario) CreateDefinitionAndInstances(instancesCount int) (err error) {

	definition := NewDefinition()
	err = c.DefinitionStorage.CreateDefinition(definition)
	if err != nil {
		log.Error().Str("type", c.StorageType).Err(err).Msg("insert document")
	}

	for i := 0; i < instancesCount; i++ {
		instance := NewInstance(definition)
		err = c.InstanceStorage.CreateInstance(instance)
		if err != nil {
			log.Error().Str("type", c.StorageType).Err(err).Msg("insert instance")
		}
	}

	return nil
}

// ListDefinitions ...
func (c *Scenario) ListDefinitions(n int) (definitions []domain.Definition, err error) {

	tenant := tenantID[RandomInt(0, 2)]
	limit := int64(n)
	offset := int64(0)

	definitions, err = c.DefinitionStorage.ListDefinitions(tenant, limit, offset)
	if err != nil {
		log.Error().Str("type", c.StorageType).Err(err).Msg("list document")
	}

	return definitions, nil
}

// ListDefinitionsAll ...
func (c *Scenario) ListDefinitionsAll(n int) (definitions []domain.Definition, err error) {

	limit := int64(n)
	offset := int64(0)

	definitions, err = c.DefinitionStorage.ListDefinitionsAll(limit, offset)
	if err != nil {
		log.Error().Str("type", c.StorageType).Err(err).Msg("list document all")
	}

	return definitions, nil
}

// GetDefinition ...
func (c *Scenario) GetDefinition(tenant string, id string) (definition *domain.Definition, err error) {

	definition, err = c.DefinitionStorage.GetDefinition(tenant, id)
	if err != nil {
		log.Error().Str("type", c.StorageType).Err(err).Msg("get definition")
	}

	return definition, nil
}

// GetDefinitionWithStat ...
func (c *Scenario) GetDefinitionWithStat(tenant string, id string) (definition *domain.Definition, err error) {

	definition, err = c.DefinitionStorage.GetDefinitionWithStats(tenant, id)
	if err != nil {
		log.Error().Str("type", c.StorageType).Err(err).Msg("get definition with stats")
	}

	return definition, nil
}

// FastWriteRead ...
func (c *Scenario) FastWriteRead() (err error) {

	definition := NewDefinition()
	err = c.DefinitionStorage.CreateDefinition(definition)
	if err != nil {
		log.Error().Str("type", c.StorageType).Err(err).Msg("insert document")
	}

	definition, err = c.DefinitionStorage.GetDefinition(definition.TenantID, definition.ID)
	if err != nil {
		log.Error().Str("type", c.StorageType).Err(err).Msg("get definition")
	}

	if definition == nil {
		if err != nil {
			log.Error().Str("type", c.StorageType).Msg("definition is nil")
		}
	}

	return nil
}

// DeleteAll ...
func (c *Scenario) DeleteAll() (err error) {

	err = c.DefinitionStorage.DeleteAll()
	if err != nil {
		log.Error().Str("type", c.StorageType).Err(err).Msg("delete definitions")
	}

	err = c.InstanceStorage.DeleteAll()
	if err != nil {
		log.Error().Str("type", c.StorageType).Err(err).Msg("delete instances")
	}

	return nil
}
