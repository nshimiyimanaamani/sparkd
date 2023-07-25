package main

import (
	"github.com/quarksgroup/sparkd/internal/db"
	"github.com/quarksgroup/sparkd/store/schema"
	lgg "github.com/sirupsen/logrus"
)

// return database instance connection
func provideDB(lg *lgg.Logger, name string) *db.DB {

	lg.Infof("connecting to database '%s'", name)

	db, err := db.New(name, lg)
	if err != nil {
		lg.Debugf("failed to connect to database: %v", err)
		lg.Fatal(err)
	}

	if ok, _ := db.IsPrimary(); ok {

		m, err := schema.Migrate(db, schema.Up)
		if err != nil {
			lg.Fatal(err)
		}

		lg.Infof("applied '%d' new migration(s)", m)
	}

	return db
}
