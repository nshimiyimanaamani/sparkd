package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/quarksgroup/sparkd/handlers/api/machines"
	"github.com/quarksgroup/sparkd/internal/core"
)

// Server is the struct that implements the http.Handler interface
type Server struct {
	vms core.MachineStore
}

func New(vms core.MachineStore) *Server {
	return &Server{vms: vms}
}

func (srv *Server) Handler() http.Handler {
	r := chi.NewRouter()

	r.Route("/machines", func(r chi.Router) {
		r.Post("/", machines.Create(srv.vms))
		r.Get("/list", machines.List(srv.vms))
		r.Route("/{vm_id}", func(r chi.Router) {
			r.Get("/", machines.Find(srv.vms))
			r.Put("/", machines.Stop())
			r.Post("/", machines.Resume())
			r.Delete("/", machines.Delete())
			r.Get("/config", machines.Config())
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("404 - Endpoint Not Found"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
	})

	return r
}
