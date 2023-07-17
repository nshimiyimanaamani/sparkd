package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/quarksgroup/sparkd/handlers/api"
	"github.com/quarksgroup/sparkd/internal/config"
	"github.com/quarksgroup/sparkd/internal/services/firecracker/vmms"
	lgg "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

var (
	// VERSION is the version of the application
	VERSION = "0.0.1"
	PREFIX  = "SPARKD"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	lg := lgg.New()

	lg.SetFormatter(&lgg.JSONFormatter{})
	lg.SetOutput(os.Stdout)
	lg.SetLevel(lgg.DebugLevel)

	cfg, err := config.Load(PREFIX)
	if err != nil {
		panic(err)
	}

	db := provideDB(lg, cfg.DBNAME)
	machines := provideMachineStore(db, cfg)
	srv := api.New(machines)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()

	go func() {
		select {
		case <-signalChan: // first signal, cancel context
			cancel()
		case <-ctx.Done():
		}
		<-signalChan // second signal, hard exit
		os.Exit(1)
	}()

	// for killing all running VMs
	defer func() {
		vmms.Cleanup()
		cancel()
	}()

	g := errgroup.Group{}

	lg.Infof("Listening on port %s", cfg.PORT)

	g.Go(func() error {
		return http.ListenAndServe(fmt.Sprintf(":%s", cfg.PORT), server(srv, lg))
	})

	<-ctx.Done()

	lg.Infoln("shutting down app")

	if err := g.Wait(); err != nil {
		lg.Fatal("main: runtime program terminated")
	}
}
