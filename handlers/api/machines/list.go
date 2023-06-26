package machines

import (
	"encoding/json"
	"net/http"

	"github.com/quarksgroup/sparkd/internal/core"
	"github.com/quarksgroup/sparkd/internal/render"
)

// For getting all running vms
func List() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		log := render.GetLogger(r.Context())

		var resp []CreateResponse = make([]CreateResponse, 0)

		for _, v := range core.RunVms {
			pid, _ := v.Vm.PID()
			resp = append(resp, CreateResponse{
				Name:   v.Name,
				State:  string(v.State),
				IpAddr: string(v.Vm.Cfg.MmdsAddress),
				ID:     v.Vm.Cfg.VMID,
				PID:    int64(pid),
			})
		}

		response, err := json.Marshal(&resp)
		if err != nil {
			log.Fatalf("failed to marshal json, %s", err)
		}
		render.JSON(w, response, http.StatusOK)
	}
}
