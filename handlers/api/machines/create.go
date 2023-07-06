package machines

import (
	"context"
	"net/http"

	"github.com/quarksgroup/sparkd/internal/core"
	"github.com/quarksgroup/sparkd/internal/render"
	"github.com/quarksgroup/sparkd/internal/services/firecracker/vmms"
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
			log.Printf("failed to create machine %v", err)
			// return
		}

		resp := CreateResponse{
			ID:         m.Id,
			SocketPath: m.SocketPath,
			Name:       in.Name,
			State:      string(m.State),
			IpAddr:     opts.FcIP,
			Agent:      m.Agent,
			CreatedAt:  m.CreatedAt,
		}

		render.JSON(w, resp, http.StatusOK)

		go (func() {
			m, err := opts.Start(r.Context(), m)
			if err != nil {
				log.Printf("failed to start created machine vm %v", err)
				// return err
			}
			core.RunVms[m.Id] = m
		})()

	}

}

func submitJob(ctx context.Context, opts *vmms.Options, m *core.Firecracker) {

	log := render.GetLogger(ctx)

	m, err := opts.Start(ctx, m)
	if err != nil {
		log.Printf("failed to start created machine vm %v", err)
		// return err
	}

	core.RunVms[m.Id] = m
}
