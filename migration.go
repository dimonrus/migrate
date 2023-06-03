package migrate

import (
	"database/sql"
	_ "embed"
	"fmt"
	"github.com/dimonrus/godb/v2"
	"os"
	"text/template"
	"time"
)

//go:embed script.tmpl
var MigrationTemplate string

// Migration struct
type Migration struct {
	// Database connection
	DBO *godb.DBO
	// Migration registry
	Registry MigrationRegistry
	// Path to migration file
	MigrationPath string
	// Registry Path in file system
	RegistryPath string
	// Registry X Path in code
	RegistryXPath string
}

// Upgrade apply migration
func (m *Migration) Upgrade(class string) error {
	for _, migration := range m.Registry[class] {
		tx, err := m.DBO.Begin()
		if err != nil {
			return err
		}
		var applyTime uint64
		// Check migration that already applied
		query := fmt.Sprintf("SELECT apply_time FROM migration_%s WHERE version = '%s'", class, migration.GetVersion())
		err = m.DBO.QueryRow(query).Scan(&applyTime)
		if err != nil && err != sql.ErrNoRows {
			_ = tx.Rollback()
			return err
		}
		// If already applied continue
		if applyTime != 0 {
			_ = tx.Commit()
			continue
		}
		// Apply new migration
		err = migration.Up(tx)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		// Update migration table
		query = fmt.Sprintf("INSERT INTO migration_%s (version, apply_time) VALUES (?, ?);", class)
		_, err = tx.Exec(query, migration.GetVersion(), time.Now().Unix())
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		err = tx.Commit()
		if err != nil {
			return err
		}
	}
	return nil
}

// Downgrade revert changes
func (m *Migration) Downgrade(class string, version string) error {
	var applyTime uint64
	for _, migration := range m.Registry[class] {
		if migration.GetVersion() != version {
			continue
		}
		// Check migration that already applied
		query := fmt.Sprintf("SELECT apply_time FROM migration_%s WHERE version = '%s'", class, version)
		err := m.DBO.QueryRow(query).Scan(&applyTime)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		// If no migration found
		if applyTime == 0 {
			m.DBO.Logger.Warnln("No migration for downgrade")
			return nil
		}
		// Begin
		tx, err := m.DBO.Begin()
		if err != nil {
			return err
		}
		// Downgrade migration
		err = migration.Down(tx)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		// Delete from migration table
		query = fmt.Sprintf("DELETE FROM migration_%s WHERE version = '%s';", class, version)
		_, err = tx.Exec(query)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		err = tx.Commit()
		if err != nil {
			return err
		}
	}
	return nil
}

// InitMigration Init migration
func (m *Migration) InitMigration(class string) error {
	var err error
	tableName := fmt.Sprintf("migration_%s", class)
	if !godb.IsTableExists(m.DBO, tableName, "") {
		query := `CREATE TABLE migration_%s (version TEXT NOT NULL, apply_time BIGINT NOT NULL);`
		_, err = m.DBO.Exec(fmt.Sprintf(query, class))
	}
	return err
}

// CreateMigrationFile Create migration file
func (m *Migration) CreateMigrationFile(class string, name string) error {
	fileName := fmt.Sprintf("m_%v_%s", time.Now().Unix(), name)
	folderPath := fmt.Sprintf("%s/%s", m.MigrationPath, class)
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return err
	}
	filePath := fmt.Sprintf("%s/%s.go", folderPath, fileName)
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	migrationTemplate := template.Must(template.New("").Parse(MigrationTemplate))
	err = migrationTemplate.Execute(f, struct {
		Class             string
		MigrationTypeName string
		RegistryPath      string
		RegistryXPath     string
	}{
		Class:             class,
		MigrationTypeName: fileName,
		RegistryPath:      m.RegistryPath,
		RegistryXPath:     m.RegistryXPath,
	})
	if err != nil {
		return err
	}
	m.DBO.Logger.Printf("Migration created: %s", filePath)
	return nil
}
