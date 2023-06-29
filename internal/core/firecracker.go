package core

import (
	"context"
	"net"

	"github.com/firecracker-microvm/firecracker-go-sdk"
)

var (
	RunVms map[string]*Firecracker = make(map[string]*Firecracker)
	IpByte byte                    = 3
)

// vmState this kind of vm-machine status
type VmState string

// avaliable vmState kind status
const (
	StateRunning VmState = "running"
	StateCreated VmState = "created"
	StateStarted VmState = "started"
	StateFailed  VmState = "failed"
	StateStopped VmState = "stopped"
)

type Firecracker struct {
	Id        string
	Name      string
	Ctx       context.Context
	CancelCtx context.CancelFunc
	Vm        *firecracker.Machine
	State     VmState
	Agent     net.IP
}
