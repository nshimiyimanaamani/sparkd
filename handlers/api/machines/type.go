package machines

import (
	"time"
)

type CreateRequest struct {
	Name  string `json:"name" validate:"required"`
	Image string `json:"image" validate:"required"`
}

type CreateResponse struct {
	ID         string     `json:"id,omitempty"`
	PID        int64      `json:"pid,omitempty"`
	SocketPath string     `json:"socket_path,omitempty"`
	State      string     `json:"state,omitempty"`
	Name       string     `json:"name,omitempty"`
	IpAddr     string     `json:"ip_address,omitempty"`
	Image      string     `json:"image,omitempty"`
	Agent      any        `json:"agent,omitempty"`
	Instance   any        `json:"instance,omitempty"`
	Resource   any        `json:"resource,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}

type DeleteRequest struct {
	ID string `json:"id" validate:"required"`
}

// Msg
type Msg struct {
	Message string `json:"message"`
}
