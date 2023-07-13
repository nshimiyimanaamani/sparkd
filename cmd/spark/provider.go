package main

import (
	"github.com/quarksgroup/sparkd/internal/config"
	"github.com/quarksgroup/sparkd/internal/core"
	"github.com/quarksgroup/sparkd/internal/db"
	"github.com/quarksgroup/sparkd/internal/services/firecracker/vmms"
	"github.com/quarksgroup/sparkd/store/machines"
)

// return mechanes interface implementation
func provideMachineStore(db *db.DB, cfg *config.Config) core.MachineStore {
	opt := vmms.New(cfg.PARENTDIR, cfg.KERNEL, cfg.FCBIN, cfg.InitrdPath, cfg.LogLevel)
	return machines.New(db, opt)
}
