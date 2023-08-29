package machines

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/quarksgroup/sparkd/internal/core"
	"github.com/quarksgroup/sparkd/internal/db"
)

// Find responsible to retrieve saved vm-machine from database using id
func (s *Store) Get(ctx context.Context, id string) (*core.Machine, error) {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	m := new(core.Machine)

	var (
		ip   sql.NullString
		path sql.NullString
	)

	if err := tx.QueryRowContext(ctx, findQuery, id).Scan(
		&m.Id,
		&m.Name,
		&m.Image,
		&m.State,
		&ip,
		&path,
		&m.CreatedAt,
		&m.UpdatedAt,
	); err != nil {
		if db.IsNoRowsError(err) {
			return nil, fmt.Errorf("machine with name %s not found", id)
		} else {
			return nil, err
		}
	}

	if ip.Valid {
		m.IpAddr = ip.String
	}

	if path.Valid {
		m.SocketPath = path.String
	}

	return m, tx.Commit()
}

var findQuery = `
SELECT
	id, name, image, state, ip_addr,socket, created_at, updated_at
FROM
	machines
WHERE
	name = $1
`
