package machines

import (
	"github.com/quarksgroup/sparkd/internal/core"
	"github.com/quarksgroup/sparkd/internal/db"
	"github.com/quarksgroup/sparkd/internal/services/firecracker/vmms"
)

type Store struct {
	db *db.DB
	// sec *db.DB // secondary database instance that hold the replace of db in case of failure
	m *vmms.Config
}

func New(db *db.DB, m *vmms.Config) *Store {
	return &Store{
		db: db,
		m:  m,
	}
}

var _ core.MachineStore = (*Store)(nil)
