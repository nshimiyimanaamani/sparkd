package machines

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/iradukunda1/firecrackerland/internal/core"
	"github.com/iradukunda1/firecrackerland/internal/render"
)

// For stopping vm using supplied vm id
func Stop() http.HandlerFunc {

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

		if err := running.Vm.PauseVM(running.Ctx); err != nil {
			log.Fatalf("failed to pause vm, %s", err)
		}
		defer running.CancelCtx()

		res, err := json.Marshal(&Msg{Message: "vm stopped successfully"})
		if err != nil {
			log.Fatalf("failed to marshal json, %s", err)
		}
		render.JSON(w, res, http.StatusOK)
	}
}
