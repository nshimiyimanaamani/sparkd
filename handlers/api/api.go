package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/quarksgroup/sparkd/handlers/api/machines"
)

func Handler() http.Handler {
	r := chi.NewRouter()

	r.Route("/machines", func(r chi.Router) {
		r.Post("/", machines.Create())
		r.Get("/list", machines.List())
		r.Post("/resume", machines.Resume())
		r.Route("/{vm_id}", func(r chi.Router) {
			r.Get("/", machines.Find())
			r.Put("/", machines.Stop())
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
