package machines

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/quarksgroup/sparkd/internal/core"
	"github.com/quarksgroup/sparkd/internal/db"
)

// List responsible to retrieve all vm instances from database
func (s *Store) List(ctx context.Context) ([]*core.Machine, error) {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	out := make([]*core.Machine, 0)
	rows, err := tx.QueryContext(ctx, listQuery)
	if err != nil {
		if db.IsNoRowsError(err) {
			return nil, fmt.Errorf("no machines found")
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		m := new(core.Machine)
		var (
			ip   sql.NullString
			path sql.NullString
		)

		if err := rows.Scan(
			&m.Id,
			&ip,
			&m.Name,
			&m.Image,
			&m.State,
			&path,
			&m.UpdatedAt,
			&m.CreatedAt,
		); err != nil {
			return nil, err
		}

		if ip.Valid {
			m.IpAddr = ip.String
		}

		if path.Valid {
			m.SocketPath = path.String
		}

		out = append(out, m)
	}

	return out, tx.Commit()
}

var listQuery = `
SELECT
	id,
	ip_addr,
	name,
	image,
	state,
	socket,
	updated_at,
	created_at
FROM
	machines
`
