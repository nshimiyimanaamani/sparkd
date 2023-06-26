package machines

import (
	"encoding/json"
	"fmt"
	"io"
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

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatalf("failed to read body, %s", err)
		}
		defer r.Body.Close()

		in := new(CreateRequest)

		if err := json.Unmarshal([]byte(body), in); err != nil {
			log.Fatalf("error during reading passed request body: %v", err.Error())
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

		// response, err := json.Marshal(&resp)
		// if err != nil {
		// 	log.Fatalf("failed to marshal json, %s", err)
		// }
		// w.Header().Add("Content-Type", "application/json")
		// w.Write(response)

		render.JSON(w, resp, http.StatusOK)

		m, err = opts.Start(m)
		if err != nil {
			fmt.Printf("failed to start vm, %s", err)
			return
		}

		core.RunVms[m.Id] = m

		return
	}

}
