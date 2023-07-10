package vmms

import (
	"context"
	"fmt"
	"time"

	"github.com/quarksgroup/sparkd/internal/core"
)

// StartVm is responsible to start vm
func Start(ctx context.Context, m *core.Firecracker) error {

	core.RunVms[m.Id] = m

	now := time.Now().UTC()
	m.UpdatedAt = &now

	if err := m.Vm.Start(ctx); err != nil {
		m.State = core.StateFailed
		return fmt.Errorf("failed to start machine: %v", err)
	}

	installSignalHandlers(ctx, m.Vm)

	m.Vm.Wait(ctx)

	m.State = core.StateRunning

	return nil
}
