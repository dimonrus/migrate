package test

import (
	"fmt"
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/godb/v2"
	"github.com/dimonrus/migrate/test/migrations"
	_ "github.com/dimonrus/migrate/test/migrations/schema"
	// _ "github.com/lib/pq"
	"testing"
)

type postgresDockerConnection struct{}

func (c *postgresDockerConnection) String() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"0.0.0.0", 5432, "migrate", "migrate", "migrate")
}

func (c *postgresDockerConnection) GetDbType() string {
	return "postgres"
}

func (c *postgresDockerConnection) GetMaxConnection() int {
	return 200
}

func (c *postgresDockerConnection) GetMaxIdleConns() int {
	return 15
}

func (c *postgresDockerConnection) GetConnMaxLifetime() int {
	return 50
}

func getPostgresDockerConnection() (*godb.DBO, error) {
	return godb.DBO{
		Options: godb.Options{
			Debug:          true,
			QueryProcessor: godb.PreparePositionalArgsQuery,
			Logger:         gocli.NewLogger(gocli.LoggerConfig{}),
		},
		Connection: &postgresDockerConnection{},
	}.Init()
}

func TestPostgresMigrations(t *testing.T) {
	q, err := getPostgresDockerConnection()
	if err != nil {
		t.Fatal(err)
	}
	migrations.Migration.DBO = q
	t.Run("init_migration", func(t *testing.T) {
		err = migrations.Migration.InitMigration("schema")
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("create_migration_file", func(t *testing.T) {
		err = migrations.Migration.CreateMigrationFile("schema", "first_migration")
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("upgrade_migration", func(t *testing.T) {
		err = migrations.Migration.Upgrade("schema")
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("downgrade_migration", func(t *testing.T) {
		err = migrations.Migration.Downgrade("schema", "name")
		if err != nil {
			t.Fatal(err)
		}
	})
}
