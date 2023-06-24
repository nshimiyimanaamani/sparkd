package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/chi"
	mddl "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/iradukunda1/firecrackerland/handlers/api"
	"github.com/iradukunda1/firecrackerland/handlers/middleware"
	"github.com/iradukunda1/firecrackerland/internal/services/vmms"
	lgg "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	lg := lgg.New()

	lg.SetFormatter(&lgg.JSONFormatter{})
	lg.SetOutput(os.Stdout)
	lg.SetLevel(lgg.DebugLevel)

	r := chi.NewMux()
	r.Use(corsHandler)
	r.Use(mddl.Recoverer)
	r.Use(middleware.SetLoggerCtx(lg))
	r.Mount("/api", api.Handler())

	lg.Infof("Listening on port 8080")

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

	g.Go(func() error {
		return http.ListenAndServe(":8080", r)
	})

	<-ctx.Done()

	lg.Infoln("shutting down app")

	if err := g.Wait(); err != nil {
		lg.Fatal("main: runtime program terminated")
	}
}

var corsHandler = cors.Handler(cors.Options{
	AllowedOrigins:   []string{"*"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: false,
	MaxAge:           300,
})
