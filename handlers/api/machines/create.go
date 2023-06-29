package machines

import (
	"fmt"
	"net/http"

	"github.com/quarksgroup/sparkd/internal/core"
	"github.com/quarksgroup/sparkd/internal/render"
	"github.com/quarksgroup/sparkd/internal/services/vmms"
)

// Create handler is for creating new vm instance
func Create() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		log := render.GetLogger(r.Context())

		IpByte := core.IpByte + 1

		core.IpByte = IpByte

		in := new(CreateRequest)

		if err := render.DecodeJSON(r, &in); err != nil {
			log.Fatalf("error during reading passed request body: %v", err.Error())
			return
		}

		opt := vmms.Options(core.Config{})

		opts, err := opt.GenerateOpt(IpByte, in.Image, in.Name)
		if err != nil {
			log.Fatalf("failed to generate option config, %s", err)
		}

		m, err := opts.Create(r.Context())
		if err != nil {
			fmt.Printf("failed to start and create vm %v", err)
			return
		}

		resp := CreateResponse{
			ID:     m.Id,
			Name:   in.Name,
			State:  string(m.State),
			IpAddr: opts.FcIP,
			Agent:  m.Agent,
		}

		render.JSON(w, resp, http.StatusOK)

		out := make(chan *core.Firecracker)

		go func() error {
			res, err := opts.Start(m)
			if err != nil {
				fmt.Printf("failed to start and create vm %v", err)
				return err
			}
			out <- res
			return nil
		}()

		m = <-out

		core.RunVms[m.Id] = m

	}

}
