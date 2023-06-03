// schema Migration file
package schema

import (
    "github.com/dimonrus/godb/v2"
    "github.com/dimonrus/migrate/test/migrations"
)

type m_1685819030_first_migration struct {}

func init() {
    migrations.Migration.Registry["schema"] = append(migrations.Migration.Registry["schema"], m_1685819030_first_migration{})
}

func (m m_1685819030_first_migration) GetVersion () string {
    return "m_1685819030_first_migration"
}

func (m m_1685819030_first_migration) Up (tx *godb.SqlTx) error {
    // write code here
    return nil
}

func (m m_1685819030_first_migration) Down (tx *godb.SqlTx) error {
    // write code here
    return nil
}