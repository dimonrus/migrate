package migrations

import "github.com/dimonrus/migrate"

var Migration = &migrate.Migration{
	RegistryPath:  "github.com/dimonrus/migrate/test/migrations",
	MigrationPath: "migrations",
	RegistryXPath: "migrations.Migration.Registry",
	Registry:      make(migrate.MigrationRegistry),
}
