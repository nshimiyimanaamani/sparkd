package machines

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/quarksgroup/sparkd/internal/core"
	"github.com/quarksgroup/sparkd/internal/render"
	"github.com/quarksgroup/sparkd/internal/services/firecracker/client"
)

// For getting vm details using supplied vm id
func Find() http.HandlerFunc {

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

		var (
			instance, resources any = nil, nil
			err                 error
		)

		if running.State != core.StateRunning {
			cli := client.NewClient(r.Context(), running.Vm.Cfg.SocketPath)
			resources, err = cli.GetResource()
			if err != nil {
				log.Fatalf("failed to get vm config, %s", err)
				render.JSON(w, err, http.StatusInternalServerError)
				return
			}

			instance, err = cli.GetInstance(r.Context())
			if err != nil {
				log.Fatalf("failed to get vm config, %s", err)
				render.JSON(w, err, http.StatusInternalServerError)
				return
			}
		}
		resp := &CreateResponse{
			Name:     running.Name,
			State:    string(running.State),
			IpAddr:   string(running.Vm.Cfg.MmdsAddress),
			ID:       running.Vm.Cfg.VMID,
			Agent:    running.Agent,
			Instance: instance,
			Resource: resources,
		}

		render.JSON(w, resp, http.StatusOK)
	}
}
