package machines

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/quarksgroup/sparkd/internal/core"
	"github.com/quarksgroup/sparkd/internal/render"
	"github.com/quarksgroup/sparkd/internal/services/firecracker/client"
)

// Config handler is for getting vm config
func Config() http.HandlerFunc {

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

		if running.State != core.StateRunning {
			res := &Msg{
				Message: fmt.Sprintf("the vm machine with this id %s is not running", id),
			}
			log.Error(res.Message)
			render.JSON(w, res, http.StatusNotFound)
			return
		}

		cli := client.NewClient(r.Context(), running.Vm.Cfg.SocketPath)

		cfg, err := cli.GetVmConfig(r.Context())
		if err != nil {
			log.Fatalf("failed to get vm config, %s", err)
			render.JSON(w, err, http.StatusInternalServerError)
			return
		}

		render.JSON(w, cfg, http.StatusOK)
	}
}
