package vmms

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	firecracker "github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/quarksgroup/sparkd/internal/core"
	llg "github.com/sirupsen/logrus"
)

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
