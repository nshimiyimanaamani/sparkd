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
		r.Get("/{vm_id}", machines.Find())
		r.Put("/{vm_id}", machines.Stop())
		r.Get("/list", machines.List())
		r.Post("/resume", machines.Resume())
		r.Delete("/{vm_id}", machines.Delete())
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("404 - Endpoint Not Found"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
	})

	return r
}
