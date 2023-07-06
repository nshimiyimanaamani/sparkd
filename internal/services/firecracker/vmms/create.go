package vmms

import (
	"context"
	"fmt"
	"time"

	firecracker "github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/quarksgroup/sparkd/internal/cmd"
	"github.com/quarksgroup/sparkd/internal/core"
	"github.com/quarksgroup/sparkd/internal/render"
	log "github.com/sirupsen/logrus"
)

// CreateVmm is responsible to create vm and return its ip address
func (o *Options) Create(ctx context.Context) (*core.Firecracker, error) {

	llg := render.GetLogger(ctx)

	cfg := o.getFcConfig()

	// logger := logging.NewFileLogger("/path/to/firecracker.log", logging.Debug)
	machineOpts := []firecracker.Opt{
		firecracker.WithLogger(log.NewEntry(llg)),
		// firecracker.WithLogger(o.Logger.WithField("app-id", o.Id)),
	}

	if err := cmd.ExposeToJail(o.RootFsImage, *cfg.JailerCfg.UID, *cfg.JailerCfg.GID); err != nil {
		return nil, fmt.Errorf("failed to expose fs to jail: %v", err)
	}

	// remove old socket path if it exists
	if _, err := cmd.RunNoneSudo(fmt.Sprintf("rm -f %s > /dev/null || true", o.ApiSocket)); err != nil {
		return nil, fmt.Errorf("failed to delete old socket path: %s", err)
	}

	if err := o.setNetwork(); err != nil {
		return nil, fmt.Errorf("failed to set network: %s", err)
	}

	m, err := firecracker.NewMachine(ctx, cfg, machineOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed creating machine: %v", err)
	}

	now := time.Now().UTC()

	res := &core.Firecracker{
		Id:         m.Cfg.VMID,
		SocketPath: m.Cfg.SocketPath,
		Ctx:        ctx,
		Name:       o.ProvidedImage,
		// cancelCtx: nil,
		Vm:        m,
		State:     core.StateCreated,
		CreatedAt: &now,
	}

	return res, nil
}
