package machines

import (
	"net/http"

	"github.com/quarksgroup/sparkd/internal/core"
	"github.com/quarksgroup/sparkd/internal/render"
)

// For getting all running vms
func List(machines core.MachineStore) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		log := render.GetLogger(r.Context())

		out := make([]CreateResponse, 0)

		res, err := machines.List(r.Context())
		if err != nil {
			e := &Msg{
				Message: err.Error(),
			}
			log.Error(err.Error())
			render.JSON(w, e, http.StatusConflict)
			return
		}

		for _, item := range res {
			// pid, _ := v.Vm.PID()
			out = append(out, CreateResponse{
				ID:         item.Id,
				Name:       item.Name,
				SocketPath: item.SocketPath,
				State:      string(item.State),
				IpAddr:     item.IpAddr,
				Image:      item.Image,
				// IpAddr:     string(item.Vm.Cfg.MmdsAddress),
				// PID:        int64(pid),
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			})
		}

		render.JSON(w, out, http.StatusOK)
	}
}
