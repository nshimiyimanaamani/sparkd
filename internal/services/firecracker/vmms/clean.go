package vmms

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	firecracker "github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/quarksgroup/sparkd/internal/cmd"
	"github.com/quarksgroup/sparkd/internal/core"
	log "github.com/sirupsen/logrus"
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

func installSignalHandlers(ctx context.Context, m *firecracker.Machine) {

	llg := log.New()

	go func() {
		// Clear some default handlers installed by the firecracker SDK:
		signal.Reset(os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

		for {
			switch s := <-c; {
			case s == syscall.SIGTERM || s == os.Interrupt:
				llg.Printf("Caught signal: %s, requesting clean shutdown", s.String())
				if err := m.Shutdown(ctx); err != nil {
					llg.Errorf("An error occurred while shutting down Firecracker VM: %v", err)
				}
			case s == syscall.SIGQUIT:
				llg.Printf("Caught signal: %s, forcing shutdown", s.String())
				if err := m.StopVMM(); err != nil {
					llg.Errorf("An error occurred while stopping Firecracker VMM: %v", err)
				}
			}
		}
	}()
}

func Cleanup() {
	for _, run := range core.RunVms {
		run.Vm.StopVMM()
	}
	cmd.RunNoneSudo("rm -f *.ext4")
}
