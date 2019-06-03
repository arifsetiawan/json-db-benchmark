package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"github.com/arifsetiawan/json-db-benchmark/domain"
)

// DefinitionStore is
type DefinitionStore struct {
	collection     *mongo.Collection
	database       *mongo.Database
	collectionName string
}

// NewDefinitionStore is
func NewDefinitionStore(database *mongo.Database, collectionName string) *DefinitionStore {
	return &DefinitionStore{
		collection:     database.Collection(collectionName),
		database:       database,
		collectionName: collectionName,
	}
}

// Drop is
func (c *DefinitionStore) Drop() (err error) {
	err = c.database.RunCommand(
		context.Background(),
		bsonx.Doc{{"drop", bsonx.String(c.collectionName)}},
	).Err()
	if err != nil {
		return err
	}

	return nil
}

// Initialize is
func (c *DefinitionStore) Initialize() (err error) {

	// create collection
	err = c.database.RunCommand(
		context.Background(),
		bsonx.Doc{{"create", bsonx.String(c.collectionName)}},
	).Err()
	if err != nil {
		return err
	}

	c.collection = c.database.Collection(c.collectionName)

	// index
	indexView := c.collection.Indexes()

	tenantIDIndex := mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "tenant_id", Value: bsonx.Int32(1)},
			{Key: "id", Value: bsonx.Int32(1)},
		},
	}
	tenantCreatedAIndex := mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "tenant_id", Value: bsonx.Int32(1)},
			{Key: "created_at", Value: bsonx.Int32(-1)},
		},
	}

	_, err = indexView.CreateMany(context.Background(),
		[]mongo.IndexModel{tenantIDIndex, tenantCreatedAIndex})
	if err != nil {
		return err
	}

	return nil
}

// CreateDefinition ...
func (c *DefinitionStore) CreateDefinition(data *domain.Definition) (err error) {

	_, err = c.collection.InsertOne(context.Background(), data)
	if err != nil {
		return err
	}

	return nil
}

// GetDefinition is
func (c *DefinitionStore) GetDefinition(tenant string, id string) (data *domain.Definition, err error) {

	d := new(domain.Definition)

	err = c.collection.FindOne(context.Background(), bson.D{{"tenant_id", tenant}, {"id", id}}).Decode(d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// GetDefinitionWithStats is
func (c *DefinitionStore) GetDefinitionWithStats(tenant string, id string) (data *domain.Definition, err error) {

	ctx := context.Background()

	/*
		queryStr := `[{
			    "$match": {
			        "tenant_id": "` + tenant + `",
			        "id": "` + id + `"
			    }
			}, {
			    "$lookup": {
			        "from": "instances",
			        "localField": "id",
			        "foreignField": "definition",
			        "as": "instances"
			    }
			}, {
			    "$unwind": {
			        "path": "$instances",
			        "preserveNullAndEmptyArrays": true
			    }
			}, {
			    "$group": {
			        "_id": "$_id",
			        "data": {
			            "$first": "$$ROOT"
			        },
			        "running_count": {
			            "$sum": {
			                "$cond": [{
			                    "$eq": ["$instances.status.running", true]
			                }, 1, 0]
			            }
			        },
			        "failed_count": {
			            "$sum": {
			                "$cond": [{
			                    "$eq": ["$instances.status.failed", true]
			                }, 1, 0]
			            }
			        },
			        "completed_count": {
			            "$sum": {
			                "$cond": [{
			                    "$eq": ["$instances.status.completed", true]
			                }, 1, 0]
			            }
			        }
			    }
			}, {
			    "$project": {
			        "id": "$data.id",
			        "reference": "$data.reference",
			        "version": "$data.version",
			        "running_count": "$running_count",
			        "failed_count": "$failed_count",
			        "completed_count": "$completed_count"
			    }
			}]`
	*/

	pipeline := mongo.Pipeline{
		{{"$match", bson.D{
			{"tenant_id", tenant},
			{"id", id},
		}}},
		{{"$lookup", bson.D{
			{"from", "instances"},
			{"localField", "id"},
			{"foreignField", "definition"},
			{"as", "instances"},
		}}},
		{{"$unwind", bson.D{
			{"path", "$instances"},
			{"preserveNullAndEmptyArrays", true},
		}}},
		{{"$group", bson.D{
			{"_id", "$_id"},
			{"data", bson.D{{"$first", "$$ROOT"}}},
			{"running_count", bson.D{
				{"$sum", bson.D{
					{"$cond", bson.A{
						bson.D{
							{"$eq", bson.A{
								"$instances.status.running",
								true,
							}},
						},
						1,
						0,
					}},
				}},
			}},
			{"failed_count", bson.D{
				{"$sum", bson.D{
					{"$cond", bson.A{
						bson.D{
							{"$eq", bson.A{
								"$instances.status.failed",
								true,
							}},
						},
						1,
						0,
					}},
				}},
			}},
			{"completed_count", bson.D{
				{"$sum", bson.D{
					{"$cond", bson.A{
						bson.D{
							{"$eq", bson.A{
								"$instances.status.completed",
								true,
							}},
						},
						1,
						0,
					}},
				}},
			}},
		}}},
		{{"$project", bson.D{
			{"id", "$data.id"},
			{"reference", "$data.reference"},
			{"version", "$data.version"},
			{"running_count", "$running_count"},
			{"failed_count", "$failed_count"},
			{"completed_count", "$completed_count"},
		}}},
	}

	cursor, err := c.collection.Aggregate(
		ctx,
		pipeline,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var items []domain.Definition
	for cursor.Next(ctx) {
		var item domain.Definition
		cursor.Decode(&item)
		items = append(items, item)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &items[0], nil

}

// ListDefinitions is
func (c *DefinitionStore) ListDefinitions(tenant string, limit int64, offset int64) (data []domain.Definition, err error) {

	ctx := context.Background()
	filter := bson.D{{"tenant_id", tenant}}
	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{"created_at", -1}})

	cursor, err := c.collection.Find(
		ctx,
		filter,
		opts,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var items []domain.Definition
	for cursor.Next(ctx) {
		var item domain.Definition
		cursor.Decode(&item)
		items = append(items, item)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// ListDefinitionsAll is
func (c *DefinitionStore) ListDefinitionsAll(limit int64, offset int64) (data []domain.Definition, err error) {

	ctx := context.Background()
	filter := bson.D{}
	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{"created_at", -1}})

	cursor, err := c.collection.Find(
		ctx,
		filter,
		opts,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var items []domain.Definition
	for cursor.Next(ctx) {
		var item domain.Definition
		cursor.Decode(&item)
		items = append(items, item)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// DeleteAll is
func (c *DefinitionStore) DeleteAll() (err error) {
	ctx := context.Background()
	_, err = c.collection.DeleteMany(
		ctx,
		bson.D{},
	)
	if err != nil {
		return err
	}

	return nil
}
