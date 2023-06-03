package test

import (
	"fmt"
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/godb/v2"
	"github.com/dimonrus/migrate/test/migrations"
	_ "github.com/dimonrus/migrate/test/migrations/schema"
	// _ "github.com/mattn/go-sqlite3"
	"testing"
)

type sqliteDockerConnection struct{}

func (c *sqliteDockerConnection) String() string {
	return fmt.Sprintf("./migrate.db")
}

func (c *sqliteDockerConnection) GetDbType() string {
	return "sqlite3"
}

func (c *sqliteDockerConnection) GetMaxConnection() int {
	return 200
}

func (c *sqliteDockerConnection) GetMaxIdleConns() int {
	return 15
}

func (c *sqliteDockerConnection) GetConnMaxLifetime() int {
	return 50
}

func getSqliteDockerConnection() (*godb.DBO, error) {
	return godb.DBO{
		Options: godb.Options{
			Debug:  true,
			Logger: gocli.NewLogger(gocli.LoggerConfig{}),
		},
		Connection: &sqliteDockerConnection{},
	}.Init()
}

func TestSqliteMigrations(t *testing.T) {
	q, err := getSqliteDockerConnection()
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
