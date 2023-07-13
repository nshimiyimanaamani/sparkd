-- +migrate Up

-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS machines (
    -- id is the unique identifier for the machine
	id VARCHAR(255) PRIMARY KEY,
    -- name is the name of the machine this can also be used as hostname
	name TEXT NOT NULL UNIQUE,
    -- ip is the ip address of the machine
    ip_addr VARCHAR(32),
    -- image is the image of the machine
    image TEXT NOT NULL,
    -- socket is the socket of the machine
    socket VARCHAR(255),
    -- state is the state of the machine
    state VARCHAR(32) NOT NULL DEFAULT 'CREATED' CHECK(state in ('CREATED', 'RUNNING','STARTED', 'FAILED','STOPPED')),
    -- created_at is the time the machine was created
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- updated_at is the time the machine was updated
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- deleted_at is the time the machine was deleted
    deleted_at DATETIME DEFAULT NULL
);
-- +migrate StatementEnd   
-- +migration StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS id_index ON machines(id);
-- +migration StatementEnd

-- +migrate Down