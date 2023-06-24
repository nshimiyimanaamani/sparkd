package vmms

import (
	"context"
	"fmt"

	"github.com/iradukunda1/firecrackerland/internal/core"
)

// StartVm is responsible to start vm
func (*Options) Start(m *core.Firecracker) (*core.Firecracker, error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := m.Vm.Start(ctx); err != nil {

		m.State = core.StateFailed

		return m, fmt.Errorf("failed to start machine: %v", err)
	}
	defer m.Vm.StopVMM()

	go func() {
		m.Vm.Wait(ctx)
	}()

	m.State = core.StateStarted
	m.CancelCtx = cancel

	return m, nil
}
