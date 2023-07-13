package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var driver = "sqlite3"

type DB struct {
	Sql *sql.DB
	dsn string
}

func New(dsn string) (*DB, error) {
	dsn = fmt.Sprintf("file:%s?cache=shared&mode=rwc&_journal_mode=WAL", dsn)

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return &DB{db, dsn}, nil
}

// BeginTx starts a transaction and returns a wrapper Tx type. This type
// provides a reference to the database and a fixed timestamp at the start of
// the transaction. The timestamp allows us to mock time during tests as well.
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {

	tx, err := db.Sql.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Return wrapper Tx that includes the transaction start time.
	return &Tx{
		Tx: tx,
		db: db,
	}, nil
}

// Primary returns the address of the primary database.
// if the current node is the primary, it returns an empty string.
func (db *DB) Primary() (string, error) {
	primaryFilename := filepath.Join(filepath.Dir(db.dsn), ".primary")

	primary, err := os.ReadFile(primaryFilename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", err
	}
	return string(primary), nil
}

func (db *DB) IsPrimary() (bool, error) {
	primary, err := db.Primary()
	if err != nil {
		return false, err
	}
	return primary == "", nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	// Close database.
	if db.Sql != nil {
		return db.Sql.Close()
	}
	return nil
}

func TxOptions(readonly bool) *sql.TxOptions {
	return &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  readonly,
	}
}
