package machines

import "net"

type CreateRequest struct {
	Name  string `json:"name" validate:"required"`
	Image string `json:"image" validate:"required"`
}

type CreateResponse struct {
	ID       string `json:"id,omitempty"`
	PID      int64  `json:"pid,omitempty"`
	State    string `json:"state,omitempty"`
	Name     string `json:"name,omitempty"`
	IpAddr   string `json:"ip_address,omitempty"`
	Agent    net.IP `json:"agent,omitempty"`
	Instance any    `json:"instance,omitempty"`
	Resource any    `json:"resource,omitempty"`
}

type DeleteRequest struct {
	ID string `json:"id" validate:"required"`
}

// Msg
type Msg struct {
	Message string `json:"message"`
}
