package core

import (
	"context"
	"time"

	"github.com/firecracker-microvm/firecracker-go-sdk"
)

var (
	RunVms map[string]*Machine = make(map[string]*Machine)
	IpByte byte                = 3
)

// vmState this kind of vm-machine status
type State string

// avaliable vmState kind status
const (
	StateRunning State = "RUNNING"
	StateCreated State = "CREATED"
	StateStarted State = "STARTED"
	StateFailed  State = "FAILED"
	StateStopped State = "STOPPED"
)

// Machine represent a vm-machine
type Machine struct {
	// Id is the unique identifier of the vm-machine
	Id string
	// Name is the name of the vm-machine provided by the user
	Name string
	// Image is the image location of the vm-machine provided by the user
	Image string
	// SocketPath is the path of the socket file of the vm-machine
	SocketPath string
	// IpAddr is the ip address of the vm-machine
	IpAddr string
	//VmIndex is the index of the vm-machine
	VmIndex   byte
	Ctx       context.Context
	CancelCtx context.CancelFunc
	Vm        *firecracker.Machine
	State     State
	Agent     any
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

// MachineStore is the interface that wraps the basic methods to manage vm-machines
type MachineStore interface {
	// Create is responsible to create a new vm-machine
	Create(context.Context, *Machine) (*Machine, error)
	// Start is responsible to start a vm-machine
	// Start(context.Context, *Machine) error
	// // Stop is responsible to stop a vm-machine
	// Stop(context.Context, *Machine) error
	// Delete is responsible to delete a vm-machine
	Delete(context.Context, string) error
	// List is responsible to list all vm-machines
	List(context.Context) ([]*Machine, error)
	// Get is responsible to get a vm-machine
	Get(context.Context, string) (*Machine, error)
	// // Update is responsible to update a vm-machine
	// Update(context.Context, *Machine) (*Machine, error)
}

// MachineService is the interface that wraps the basic methods to manage vm-machines
type MachineService interface {
	// Create is responsible to create a new vm-machine
	Create(context.Context, *Machine) (*Machine, error)
}
