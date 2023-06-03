package migrate

import (
	"github.com/dimonrus/godb/v2"
)

// IMigrationFile migration file interface
type IMigrationFile interface {
	// Up Apply migration
	Up(tx *godb.SqlTx) error
	// Down Revert migration
	Down(tx *godb.SqlTx) error
	// GetVersion get migration version
	GetVersion() string
}

// MigrationRegistry migration registry
type MigrationRegistry map[string][]IMigrationFile
