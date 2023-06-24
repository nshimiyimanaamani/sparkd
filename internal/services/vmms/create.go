package vmms

import (
	"context"
	"fmt"

	firecracker "github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/iradukunda1/firecrackerland/internal/cmd"
	"github.com/iradukunda1/firecrackerland/internal/core"
	"github.com/iradukunda1/firecrackerland/internal/render"
	log "github.com/sirupsen/logrus"
)

// CreateVmm is responsible to create vm and return its ip address
func (o *Options) Create(ctx context.Context) (*core.Firecracker, error) {

	llg := render.GetLogger(ctx)

	cfg := o.getFcConfig()

	machineOpts := []firecracker.Opt{
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

	m, err := firecracker.NewMachine(ctx, cfg, machineOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed creating machine: %v", err)
	}

	installSignalHandlers(ctx, m)

	res := &core.Firecracker{
		Id:   m.Cfg.VMID,
		Ctx:  ctx,
		Name: o.ProvidedImage,
		// cancelCtx: nil,
		Vm:    m,
		State: core.StateCreated,
	}

	return res, nil
}
