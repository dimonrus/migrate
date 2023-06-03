package test

import (
	"fmt"
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/godb/v2"
	"github.com/dimonrus/migrate/test/migrations"
	_ "github.com/dimonrus/migrate/test/migrations/schema"
	// _ "github.com/go-sql-driver/mysql"
	"testing"
)

type mysqlDockerConnection struct{}

func (c *mysqlDockerConnection) String() string {
	//username:password@protocol(address)/dbname?param=value
	return fmt.Sprintf("root:migrate@/migrate")
}

func (c *mysqlDockerConnection) GetDbType() string {
	return "mysql"
}

func (c *mysqlDockerConnection) GetMaxConnection() int {
	return 200
}

func (c *mysqlDockerConnection) GetMaxIdleConns() int {
	return 15
}

func (c *mysqlDockerConnection) GetConnMaxLifetime() int {
	return 50
}

func getMysqlDockerConnection() (*godb.DBO, error) {
	return godb.DBO{
		Options: godb.Options{
			Debug:  true,
			Logger: gocli.NewLogger(gocli.LoggerConfig{}),
		},
		Connection: &mysqlDockerConnection{},
	}.Init()
}

func TestMysqlMigrations(t *testing.T) {
	q, err := getMysqlDockerConnection()
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
