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
	opts := []firecracker.Opt{
		// firecracker.WithProcessRunner(o.jailerCommand(ctx, cfg.VMID, true)),
		firecracker.WithLogger(log.NewEntry(llg)),
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

	m, err := firecracker.NewMachine(ctx, cfg, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create new machine instance: %w", err)
	}

	now := time.Now().UTC()

	res := &core.Firecracker{
		Id:         m.Cfg.VMID,
		SocketPath: m.Cfg.SocketPath,
		Image:      o.ProvidedImage,
		Name:       o.Name,
		Ctx:        ctx,
		// CancelCtx:  vmCancel,
		Vm:        m,
		Agent:     m.Cfg.NetworkInterfaces[0].StaticConfiguration,
		State:     core.StateCreated,
		CreatedAt: &now,
	}

	defer Start(ctx, res)

	return res, nil
}
