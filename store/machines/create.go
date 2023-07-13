package machines

import (
	"context"
	"fmt"

	"github.com/quarksgroup/sparkd/internal/core"
	"github.com/quarksgroup/sparkd/internal/db"
	"github.com/quarksgroup/sparkd/internal/render"
)

// Create is responsible to store and create new vm-machine
func (s *Store) Create(ctx context.Context, m *core.Machine) (*core.Machine, error) {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := tx.QueryRowContext(
		ctx,
		insertQuery,
		m.Id,
		m.Name,
		m.Image,
		m.State,
	).Scan(&m.CreatedAt); err != nil {
		if db.IsUniqueConstraintError(err) {
			return nil, fmt.Errorf("machine with this name already exists")
		} else {
			return nil, err
		}
	}

	go func() {

		log := render.GetLogger(ctx)

		m, err = s.m.Create(ctx, m)
		if err != nil {
			log.Errorf("failed to create new vm-machine: %v", err)
			return
		}

		_, err = tx.ExecContext(ctx, updateQuery, m.State, m.IpAddr, m.SocketPath, m.UpdatedAt, m.Id)
		if err != nil {
			log.Errorf("failed to update vm-machine: %v", err)
			return
		}

	}()

	return m, tx.Commit()
}

var insertQuery = `
INSERT INTO machines 
	(id, name, image, state)
VALUES 
	($1, $2, $3, $4) 
RETURNING created_at`

var updateQuery = `
UPDATE machines
SET
	state = $1,
	ip_addr = $2,
	socket = $3,
	updated_at = $4
WHERE
	id = $5
`
