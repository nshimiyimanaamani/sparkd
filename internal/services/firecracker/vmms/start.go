package vmms

import (
	"context"
	"fmt"
	"time"

	"github.com/quarksgroup/sparkd/internal/core"
)

// StartVm is responsible to start vm
func Start(ctx context.Context, m *core.Firecracker) error {

	vmCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	core.RunVms[m.Id] = m

	now := time.Now().UTC()
	m.UpdatedAt = &now

	if err := m.Vm.Start(vmCtx); err != nil {
		m.State = core.StateFailed
		return fmt.Errorf("failed to start machine: %v", err)
	}

	installSignalHandlers(vmCtx, m.Vm)

	// go m.Vm.Wait(vmCtx)

	go func() {
		defer m.CancelCtx()
		m.Vm.Wait(vmCtx)
	}()

	m.State = core.StateRunning
	m.CancelCtx = cancel

	return nil
}
