package vmms

import (
	"context"
	"fmt"

	"github.com/quarksgroup/sparkd/internal/core"
)

// StartVm is responsible to start vm
func (*Options) Start(m *core.Firecracker) (*core.Firecracker, error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := m.Vm.Logger()

	if err := m.Vm.Start(ctx); err != nil {

		m.State = core.StateFailed

		return m, fmt.Errorf("failed to start machine: %v", err)
	}
	defer func() {
		if err := m.Vm.StopVMM(); err != nil {
			log.Errorf("An error occurred while stopping Firecracker VMM: %v", err)
		}
	}()

	installSignalHandlers(ctx, m.Vm)

	go func() {
		defer m.CancelCtx()
		m.Vm.Wait(context.Background())
	}()

	// if err := m.Vm.Wait(context.Background()); err != nil {
	// 	return nil, fmt.Errorf("wait returned an error %s", err)
	// }

	log.Printf("Start machine was happy")

	m.State = core.StateStarted
	m.CancelCtx = cancel

	return m, nil
}
