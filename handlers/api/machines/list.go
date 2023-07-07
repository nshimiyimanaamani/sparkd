package machines

import (
	"net/http"

	"github.com/quarksgroup/sparkd/internal/core"
	"github.com/quarksgroup/sparkd/internal/render"
)

// For getting all running vms
func List() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// log := render.GetLogger(r.Context())

		out := make([]CreateResponse, 0)

		for _, v := range core.RunVms {
			pid, _ := v.Vm.PID()
			out = append(out, CreateResponse{
				Name:       v.Name,
				SocketPath: v.SocketPath,
				State:      string(v.State),
				IpAddr:     string(v.Vm.Cfg.MmdsAddress),
				ID:         v.Vm.Cfg.VMID,
				PID:        int64(pid),
				CreatedAt:  v.CreatedAt,
				UpdatedAt:  v.UpdatedAt,
			})
		}

		render.JSON(w, out, http.StatusOK)
	}
}
