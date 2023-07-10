package machines

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/quarksgroup/sparkd/internal/core"
	"github.com/quarksgroup/sparkd/internal/render"
)

// For resuming vm using supplied vm id
func Resume() http.HandlerFunc {

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

		if err := running.Vm.ResumeVM(running.Ctx); err != nil {
			log.Fatalf("failed to resume vm, %s", err)
		}

		res, err := json.Marshal(&Msg{Message: "vm resumed successfully"})
		if err != nil {
			log.Fatalf("failed to marshal json, %s", err)
		}

		render.JSON(w, res, http.StatusOK)
	}
}
