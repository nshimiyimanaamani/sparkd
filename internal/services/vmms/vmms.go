package vmms

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	firecracker "github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/quarksgroup/sparkd/internal/core"
	llg "github.com/sirupsen/logrus"
)

// StoppedOK is the VMM stopped status.
type StoppedOK = bool

var (
	// StoppedGracefully indicates the machine was stopped gracefully.
	StoppedGracefully = StoppedOK(true)
	// StoppedForcefully indicates that the machine did not stop gracefully
	// and the shutdown had to be forced.
	StoppedForcefully = StoppedOK(false)
)

func InstallSignalHandlers(ctx context.Context, m *core.Firecracker) chan bool {
	chanStopped := make(chan bool, 1)
	log := llg.New()
	go func() {
		// Clear selected default handlers installed by the firecracker SDK:
		signal.Reset(os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		fmt.Println("Caught SIGTERM, requesting clean reset")
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		fmt.Println("Caught SIGTERM, requesting clean notification")
		for {
			dfc := defaultFc{
				vm: m,
			}

			switch s := <-c; {
			case s == syscall.SIGTERM || s == os.Interrupt:
				log.Info("Caught SIGINT, requesting clean shutdown")
				chanStopped <- dfc.Stop(ctx)
			}
		}
	}()
	return chanStopped
}

func installSignalHandlers(ctx context.Context, m *firecracker.Machine) {

	log := llg.New()

	go func() {
		// Clear some default handlers installed by the firecracker SDK:
		signal.Reset(os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

		for {
			switch s := <-c; {
			case s == syscall.SIGTERM || s == os.Interrupt:
				fmt.Println("Caught SIGTERM, requesting clean shutdown")
				if err := m.Shutdown(ctx); err != nil {
					log.Errorf("Machine shutdown failed with error: %v", err)
				}
				time.Sleep(20 * time.Second)

				// There's no direct way of checking if a VM is running, so we test if we can send it another shutdown
				// request. If that fails, the VM is still running and we need to kill it.
				if err := m.Shutdown(ctx); err == nil {
					fmt.Println("Timeout exceeded, forcing shutdown") // TODO: Proper logging
					if err := m.StopVMM(); err != nil {
						log.Errorf("VMM stop failed with error: %v", err)
					}
				}
			case s == syscall.SIGQUIT:
				fmt.Println("Caught SIGQUIT, forcing shutdown")
				if err := m.StopVMM(); err != nil {
					log.Errorf("VMM stop failed with error: %v", err)
				}
			}
		}
	}()
}

func Cleanup() {
	for _, run := range core.RunVms {
		run.Vm.StopVMM()
	}
}

type defaultFc struct {
	sync.Mutex

	vm         *core.Firecracker
	wasStopped bool
}

func (m *defaultFc) Stop(ctx context.Context) StoppedOK {

	m.Lock()
	defer m.Unlock()

	log := llg.New()

	if m.vm.State != core.StateRunning {
		m.wasStopped = true
	} else {
		return StoppedGracefully
	}

	shutdownCtx, cancelFunc := context.WithTimeout(ctx, time.Second*time.Duration(30))
	defer cancelFunc()

	log.Info("Attempting VMM graceful shutdown...")

	chanStopped := make(chan error, 1)
	go func() {
		// Ask the machine to shut down so the file system gets flushed
		// and out changes are written to disk.
		chanStopped <- m.vm.Vm.Shutdown(shutdownCtx)
	}()

	stoppedState := StoppedForcefully

	select {
	case stopErr := <-chanStopped:
		if stopErr != nil {
			log.Warn("VMM stopped with error but within timeout", "reason", stopErr)
			log.Warn("VMM stopped forcefully", "error", m.vm.Vm.StopVMM())
		} else {
			log.Warn("VMM stopped gracefully")
			stoppedState = StoppedGracefully
		}
	case <-shutdownCtx.Done():
		log.Warn("VMM failed to stop gracefully: timeout reached")
		log.Warn("VMM stopped forcefully", "error", m.vm.Vm.StopVMM())
	}

	log.Info("Cleaning up CNI network...")

	// cniCleanupErr := m.cleanupCNINetwork()

	// log.Info("CNI network cleanup status", "error", cniCleanupErr)

	return stoppedState
}
