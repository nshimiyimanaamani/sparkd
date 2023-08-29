package machines

import (
	"net/http"

	"github.com/quarksgroup/sparkd/internal/core"
	"github.com/quarksgroup/sparkd/internal/rand"
	"github.com/quarksgroup/sparkd/internal/render"
)

// Create handler is for creating new vm instance
func Create(machines core.MachineStore) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		log := render.GetLogger(r.Context())

		IpByte := core.IpByte + 1

		core.IpByte = IpByte

		in := new(CreateRequest)

		if err := render.DecodeJSON(r, &in); err != nil {
			log.Errorf("error during reading passed request body: %v", err.Error())
			msg := &Msg{
				Message: err.Error(),
			}
			render.JSON(w, msg, http.StatusBadRequest)
			return
		}

		if !core.MatchName(in.Name) {
			err := &Msg{
				Message: "name must be alphanumeric",
			}
			log.Errorln(err)
			render.JSON(w, err, http.StatusBadRequest)
			return
		}

		m := &core.Machine{
			Id:      rand.UUID(),
			Name:    in.Name,
			Image:   in.Image,
			VmIndex: IpByte,
			State:   core.StateCreated,
		}

		res, err := machines.Create(r.Context(), m)
		if err != nil {
			log.Error(err.Error())
			msg := &Msg{
				Message: err.Error(),
			}
			render.JSON(w, msg, http.StatusConflict)
			return
		}

		out := &CreateResponse{
			ID:         m.Id,
			SocketPath: m.SocketPath,
			Name:       in.Name,
			State:      string(m.State),
			IpAddr:     res.IpAddr,
			Agent:      m.Agent,
		}

		render.JSON(w, out, http.StatusOK)

	}

}
