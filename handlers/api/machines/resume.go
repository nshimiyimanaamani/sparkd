package machines

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/quarksgroup/sparkd/internal/core"
	"github.com/quarksgroup/sparkd/internal/render"
)

// For resuming vm using supplied vm id
func Resume() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		log := render.GetLogger(r.Context())

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatalf("failed to read body, %s", err)
		}
		defer r.Body.Close()

		in := new(DeleteRequest)

		json.Unmarshal([]byte(body), in)
		if err != nil {
			log.Fatalf("error during reading passed request body: %v", err.Error())
		}

		running := core.RunVms[in.ID]

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
