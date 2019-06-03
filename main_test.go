package main

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/arifsetiawan/json-db-benchmark/domain"
	"github.com/arifsetiawan/json-db-benchmark/scenario"
	"github.com/arifsetiawan/json-db-benchmark/storage/arangodb"
	"github.com/arifsetiawan/json-db-benchmark/storage/mongodb"
	"github.com/arifsetiawan/json-db-benchmark/storage/mysql"
	"github.com/arifsetiawan/json-db-benchmark/storage/postgres"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

var (
	mongoScenario      *scenario.Scenario
	arangoRockScenario *scenario.Scenario
	arangoMMScenario   *scenario.Scenario
	couchbaseScenario  *scenario.Scenario
	postgresScenario   *scenario.Scenario
	postgres11Scenario *scenario.Scenario
	mysqlScenario      *scenario.Scenario

	err              error
	mongoDefinitions []domain.Definition
)

func init() {

	{
		// create client
		mongoClient, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URL")))
		if err != nil {
			log.Error().Err(err).Msg("Failed to create database connection")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		err = mongoClient.Connect(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create database connection")
		}

		// create scenario
		mongoScenario = &scenario.Scenario{
			StorageType:       "mongodb",
			DefinitionStorage: mongodb.NewDefinitionStore(mongoClient.Database(os.Getenv("MONGODB_DATABASE")), "definitions"),
			InstanceStorage:   mongodb.NewInstanceStore(mongoClient.Database(os.Getenv("MONGODB_DATABASE")), "instances"),
		}
	}

	{
		arangoConnection, err := http.NewConnection(http.ConnectionConfig{
			Endpoints: []string{os.Getenv("ARANGODB_ROCK_URL")},
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to create http connection")
		}

		arangoClient, err := driver.NewClient(driver.ClientConfig{
			Connection:     arangoConnection,
			Authentication: driver.BasicAuthentication(os.Getenv("ARANGODB_USERNAME"), os.Getenv("ARANGODB_PASSWORD")),
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to create client connection")
		}

		arangoDB, err := arangoClient.Database(nil, os.Getenv("ARANGODB_DATABASE"))
		if err != nil {
			log.Error().Err(err).Msg("Failed to open database")
		}

		arangoRockScenario = &scenario.Scenario{
			StorageType:       "arangodb-rock",
			DefinitionStorage: arangodb.NewDefinitionStore(arangoDB, "definitions"),
			InstanceStorage:   arangodb.NewInstanceStore(arangoDB, "instances"),
		}
	}

	{
		arangoConnection, err := http.NewConnection(http.ConnectionConfig{
			Endpoints: []string{os.Getenv("ARANGODB_MM_URL")},
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to create http connection")
		}

		arangoClient, err := driver.NewClient(driver.ClientConfig{
			Connection:     arangoConnection,
			Authentication: driver.BasicAuthentication(os.Getenv("ARANGODB_USERNAME"), os.Getenv("ARANGODB_PASSWORD")),
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to create client connection")
		}

		arangoDB, err := arangoClient.Database(nil, os.Getenv("ARANGODB_DATABASE"))
		if err != nil {
			log.Error().Err(err).Msg("Failed to open database")
		}

		arangoMMScenario = &scenario.Scenario{
			StorageType:       "arangodb-mm",
			DefinitionStorage: arangodb.NewDefinitionStore(arangoDB, "definitions"),
			InstanceStorage:   arangodb.NewInstanceStore(arangoDB, "instances"),
		}
	}

	{
		dbctx, dbcancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer dbcancel()

		db, err := sqlx.ConnectContext(dbctx, "postgres", os.Getenv("POSTGRES_CONNECTION_STR"))

		if dbctx.Err() != nil {
			log.Error().Err(err).Msg("Postgres connection timed out")
		}
		if err != nil {
			log.Error().Err(err).Msg("Postgres connection failed")
		}

		// create scenario
		postgresScenario = &scenario.Scenario{
			StorageType:       "postgres",
			DefinitionStorage: postgres.NewDefinitionStore(db),
			InstanceStorage:   postgres.NewInstanceStore(db),
		}
	}

	{
		dbctx, dbcancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer dbcancel()

		db, err := sqlx.ConnectContext(dbctx, "postgres", os.Getenv("POSTGRES11_CONNECTION_STR"))

		if dbctx.Err() != nil {
			log.Error().Err(err).Msg("Postgres connection timed out")
		}
		if err != nil {
			log.Error().Err(err).Msg("Postgres connection failed")
		}

		// create scenario
		postgres11Scenario = &scenario.Scenario{
			StorageType:       "postgres11",
			DefinitionStorage: postgres.NewDefinitionStore(db),
			InstanceStorage:   postgres.NewInstanceStore(db),
		}
	}

	{
		dbctx, dbcancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer dbcancel()

		db, err := sqlx.ConnectContext(dbctx, "mysql", os.Getenv("MYSQL_CONNECTION_STR"))

		if dbctx.Err() != nil {
			log.Error().Err(err).Msg("mysql connection timed out")
		}
		if err != nil {
			log.Error().Err(err).Msg("mysql connection failed")
		}

		// create scenario
		mysqlScenario = &scenario.Scenario{
			StorageType:       "mysql",
			DefinitionStorage: mysql.NewDefinitionStore(db),
			InstanceStorage:   mysql.NewInstanceStore(db),
		}
	}

	/*
		{
			cluster, err := gocb.Connect(os.Getenv("COUCHBASE_URL"))
			if err != nil {
				log.Error().Err(err).Msg("Failed to create database connection")
			}

			cluster.Authenticate(gocb.PasswordAuthenticator{
				Username: os.Getenv("COUCHBASE_ADMIN_USERNAME"),
				Password: os.Getenv("COUCHBASE_ADMIN_PASSWORD"),
			})

			bucket, err := cluster.OpenBucket(os.Getenv("COUCHBASE_BUCKET"), os.Getenv("COUCHBASE_BUCKET_PASSWORD"))
			if err != nil {
				log.Error().Err(err).Msg("Failed to connect to bucket")
			}

			couchbaseScenario = &scenario.Scenario{
				StorageType:       "couchbase",
				DefinitionStorage: couchbase.NewDefinitionStore(bucket),
				InstanceStorage:   couchbase.NewInstanceStore(bucket),
			}
		}
	*/

}

func reinitStorage(scen *scenario.Scenario) {
	scen.Drop()
	time.Sleep(2 * time.Second)
	scen.Initialize()
}

/*
func TestDropInit(t *testing.T) {
	scenarios := []struct {
		name string
		scen *scenario.Scenario
	}{
		{"mongodb", mongoScenario},
		{"arango-rock", arangoRockScenario},
		{"arango-mm", arangoMMScenario},
		{"postgres", postgresScenario},
		{"postgres11", postgres11Scenario},
		{"mysql", mysqlScenario},
		{"couchbase", couchbaseScenario},
	}

	for _, s := range scenarios {
		t.Logf("run %s", s.name)
		reinitStorage(s.scen)

		for i := 0; i < 2; i++ {
			s.scen.CreateSingleDefinition()
		}

		definitions, _ := s.scen.ListDefinitionsAll(2)

		if len(definitions) != 2 {
			t.Errorf("%s, expected length %d, got %d", s.name, 2, len(definitions))
		}
	}
}
*/

func BenchmarkInsert1(b *testing.B) {
	scenarios := []struct {
		name string
		scen *scenario.Scenario
	}{
		{"mongodb", mongoScenario},
		{"arango-rock", arangoRockScenario},
		{"arango-mm", arangoMMScenario},
		{"postgres", postgresScenario},
		{"postgres11", postgres11Scenario},
		{"mysql", mysqlScenario},
		//{"couchbase", couchbaseScenario},
	}

	for _, s := range scenarios {
		reinitStorage(s.scen)

		b.Run(s.name, func(b *testing.B) {

			for i := 0; i < b.N; i++ {
				s.scen.CreateSingleDefinition()
			}
		})
	}
}

func BenchmarkInsert100(b *testing.B) {
	scenarios := []struct {
		name string
		scen *scenario.Scenario
	}{
		{"mongodb", mongoScenario},
		{"arango-rock", arangoRockScenario},
		{"arango-mm", arangoMMScenario},
		{"postgres", postgresScenario},
		{"postgres11", postgres11Scenario},
		{"mysql", mysqlScenario},
		//{"couchbase", couchbaseScenario},
	}

	for _, s := range scenarios {
		reinitStorage(s.scen)

		b.Run(s.name, func(b *testing.B) {

			for i := 0; i < b.N; i++ {
				for i := 0; i < 100; i++ {
					s.scen.CreateSingleDefinition()
				}
			}
		})
	}
}

func BenchmarkGet(b *testing.B) {
	scenarios := []struct {
		name string
		scen *scenario.Scenario
	}{
		{"mongodb", mongoScenario},
		{"arango-rock", arangoRockScenario},
		{"arango-mm", arangoMMScenario},
		{"postgres", postgresScenario},
		{"postgres11", postgres11Scenario},
		{"mysql", mysqlScenario},
		//{"couchbase", couchbaseScenario},
	}

	for _, s := range scenarios {

		for _, n := range []int{100, 500, 1000} {
			reinitStorage(s.scen)

			for i := 0; i < n; i++ {
				s.scen.CreateSingleDefinition()
			}

			time.Sleep(2 * time.Second)

			definitions, _ := s.scen.ListDefinitions(n)

			b.Run(fmt.Sprintf("%s-%d", s.name, n), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					def := definitions[scenario.RandomInt(0, len(definitions))]
					s.scen.GetDefinition(def.TenantID, def.ID)
				}
			})
		}
	}
}

func BenchmarkList(b *testing.B) {
	scenarios := []struct {
		name string
		scen *scenario.Scenario
	}{
		{"mongodb", mongoScenario},
		{"arango-rock", arangoRockScenario},
		{"arango-mm", arangoMMScenario},
		{"postgres", postgresScenario},
		{"postgres11", postgres11Scenario},
		{"mysql", mysqlScenario},
		//{"couchbase", couchbaseScenario},
	}

	for _, s := range scenarios {

		for _, n := range []int{100, 500, 1000} {
			reinitStorage(s.scen)

			for i := 0; i < n; i++ {
				s.scen.CreateSingleDefinition()
			}

			b.Run(fmt.Sprintf("%s-%d", s.name, n), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					s.scen.ListDefinitions(100)
				}
			})
		}
	}
}

func BenchmarkInsertDI(b *testing.B) {
	scenarios := []struct {
		name string
		scen *scenario.Scenario
	}{
		{"mongodb", mongoScenario},
		{"arango-rock", arangoRockScenario},
		{"arango-mm", arangoMMScenario},
		{"postgres", postgresScenario},
		{"postgres11", postgres11Scenario},
		{"mysql", mysqlScenario},
		//{"couchbase", couchbaseScenario},
	}

	for _, s := range scenarios {
		reinitStorage(s.scen)

		b.Run(s.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s.scen.CreateDefinitionAndInstances(100)
			}
		})
	}
}

func BenchmarkGetStats(b *testing.B) {
	scenarios := []struct {
		name string
		scen *scenario.Scenario
	}{
		{"mongodb", mongoScenario},
		{"arango-rock", arangoRockScenario},
		{"arango-mm", arangoMMScenario},
		{"postgres", postgresScenario},
		{"postgres11", postgres11Scenario},
		{"mysql", mysqlScenario},
		//{"couchbase", couchbaseScenario},
	}

	for _, s := range scenarios {

		for _, n := range []int{100, 500} {
			reinitStorage(s.scen)

			for i := 0; i < n; i++ {
				s.scen.CreateDefinitionAndInstances(100)
			}

			definitions, _ := s.scen.ListDefinitions(n)

			b.Run(fmt.Sprintf("%s-%d", s.name, n), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					def := definitions[scenario.RandomInt(0, len(definitions))]
					s.scen.GetDefinitionWithStat(def.TenantID, def.ID)
				}
			})
		}
	}
}
