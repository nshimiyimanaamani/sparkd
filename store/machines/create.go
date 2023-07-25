package machines

import (
	"context"
	"fmt"
	"time"

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

	// run background process to start vm with using goroutine and channel to handle context canceled.

	cntx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	log := render.GetLogger(cntx)

	resultChan := make(chan error)
	go func() {
		_, err := s.m.Create(cntx, m)
		if err != nil {
			log.Errorf("failed to create new VM-machine: %v", err)
		}
		resultChan <- err
	}()

	select {
	case <-cntx.Done():
		fmt.Println("VM startup timed out.")
		_, err = tx.Exec(updateQuery, m.State, m.IpAddr, m.SocketPath, m.UpdatedAt, m.Id)
		if err != nil {
			log.Errorf("failed to update VM-machine: %v", err)
		}
	case err := <-resultChan:
		if err != nil {
			fmt.Println("VM creation action failed:", err)
			_, err = tx.Exec(updateQuery, m.State, m.IpAddr, m.SocketPath, m.UpdatedAt, m.Id)
			if err != nil {
				log.Errorf("failed to update VM-machine: %v", err)
			}
		} else {
			fmt.Println("VM startup completed successfully.")
			_, err = tx.Exec(updateQuery, m.State, m.IpAddr, m.SocketPath, m.UpdatedAt, m.Id)
			if err != nil {
				log.Errorf("failed to update VM-machine: %v", err)
			}
		}
	}

	m.CancelCtx = cancel

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
