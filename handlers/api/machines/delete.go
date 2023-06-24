package machines

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/iradukunda1/firecrackerland/internal/core"
	"github.com/iradukunda1/firecrackerland/internal/render"
)

// for deleting supplied vm id
func Delete() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		log := render.GetLogger(r.Context())

		id := chi.URLParam(r, "vm_id")

		running, ok := core.RunVms[id]
		if !ok {
			res := &Msg{
				Message: fmt.Sprintf("the vm machine with this id %s is not exist", id),
			}
			log.Error(res.Message)
			render.JSON(w, res, http.StatusNotFound)
			return
		}

		if err := running.Vm.Shutdown(running.Ctx); err != nil {
			log.Fatalf("failed to delete vm, %s", err)
			render.JSON(w, err, http.StatusInternalServerError)
			return
		}

		res, err := json.Marshal(&Msg{Message: "vm deleted successfully"})
		if err != nil {
			log.Fatalf("failed to marshal json, %s", err)
			render.JSON(w, err, http.StatusInternalServerError)
			return
		}

		delete(core.RunVms, id)

		render.JSON(w, res, http.StatusOK)
	}
}
