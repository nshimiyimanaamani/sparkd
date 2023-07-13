package main

import (
	"net/http"

	"github.com/go-chi/chi"
	mddl "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/quarksgroup/sparkd/handlers/api"
	"github.com/quarksgroup/sparkd/handlers/middleware"
	log "github.com/sirupsen/logrus"
)

func server(srv *api.Server, lg *log.Logger) http.Handler {
	r := chi.NewMux()
	r.Use(corsHandler)
	r.Use(mddl.Recoverer)
	r.Use(middleware.SetLoggerCtx(lg))
	r.Mount("/api", srv.Handler())
	return r
}

var corsHandler = cors.Handler(cors.Options{
	AllowedOrigins:   []string{"*"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: false,
	MaxAge:           300,
})
