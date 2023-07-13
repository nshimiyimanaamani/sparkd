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
func (dfOpt *Config) Create(ctx context.Context, fc *core.Machine) (*core.Machine, error) {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	llg := render.GetLogger(ctx)

	opt, err := dfOpt.generateOpt(fc.VmIndex, fc.Image, fc.Id, fc.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to generate option config, %s", err)
	}

	fcCfg := getFcConfig(opt)

	opts := []firecracker.Opt{
		firecracker.WithLogger(log.NewEntry(llg)),
	}

	if err := cmd.ExposeToJail(opt.rootFsImage, *fcCfg.JailerCfg.UID, *fcCfg.JailerCfg.GID); err != nil {
		return nil, fmt.Errorf("failed to expose fs to jail: %v", err)
	}

	// remove old socket path if it exists
	if _, err := cmd.RunNoneSudo(fmt.Sprintf("rm -f %s > /dev/null", opt.apiSocket)); err != nil {
		return nil, fmt.Errorf("failed to delete old socket path: %s", err)
	}

	if err := opt.setNetwork(); err != nil {
		return nil, fmt.Errorf("failed to set network: %s", err)
	}

	m, err := firecracker.NewMachine(ctx, fcCfg, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create new machine instance: %w", err)
	}

	now := time.Now().UTC()

	fc.SocketPath = m.Cfg.SocketPath
	fc.Ctx = ctx
	fc.Vm = m
	fc.CancelCtx = cancel
	fc.Agent = m.Cfg.NetworkInterfaces[0].StaticConfiguration
	fc.CreatedAt = &now

	defer start(ctx, fc)

	return fc, nil
}
