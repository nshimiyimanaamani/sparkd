package vmms

import (
	"context"
	"fmt"
	"time"

	"github.com/quarksgroup/sparkd/internal/core"
)

// StartVm is responsible to start vm
func (*Options) Start(ctx context.Context, m *core.Firecracker) (*core.Firecracker, error) {

	ctx, cancel := context.WithTimeout(ctx, 500*time.Second)
	defer cancel()

	log := m.Vm.Logger()
	now := time.Now().UTC()

	m.UpdatedAt = &now

	if err := m.Vm.Start(context.Background()); err != nil {

		m.State = core.StateFailed
		return m, fmt.Errorf("failed to start machine: %v", err)
	}

	installSignalHandlers(ctx, m.Vm)

	// go func() {
	// 	defer m.CancelCtx()
	if err := m.Vm.Wait(ctx); err != nil {
		return nil, fmt.Errorf("wait returned an error %s", err)
	}
	// }()

	// if err := m.Vm.Wait(context.Background()); err != nil {
	// 	return nil, fmt.Errorf("wait returned an error %s", err)
	// }

	log.Printf("Start machine was happy")

	m.State = core.StateStarted
	m.CancelCtx = cancel

	return m, nil
}
