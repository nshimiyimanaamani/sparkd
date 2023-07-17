package schema

import (
	"embed"
	"fmt"

	"github.com/quarksgroup/sparkd/internal/db"
	migrate "github.com/rubenv/sql-migrate"
)

const root = "migrations"

//go:embed migrations/*
var migrations embed.FS

// Direction ...
type Direction migrate.MigrationDirection

// Migration directions
const (
	// Migration apply
	Up Direction = 0
	// Migration Rollback
	Down Direction = 1
)

// Migrate peforms database migrations and returns an error
// if migration fails.
func Migrate(db *db.DB, dir Direction) (int, error) {
	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: migrations,
		Root:       root,
	}

	n, err := migrate.Exec(db.Sql, "sqlite3", migrations, migrate.MigrationDirection(dir))
	if err != nil {
		return n, fmt.Errorf("could not apply migrations%w", err)
	}
	return n, nil
}
