package client

import (
	"context"
	"fmt"

	"github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/firecracker-microvm/firecracker-go-sdk/client/models"
	"github.com/quarksgroup/sparkd/internal/render"
	"github.com/sirupsen/logrus"
)

type Client struct {
	*firecracker.Client
}

func NewClient(ctx context.Context, socketPath string) *Client {

	logs := render.GetLogger(ctx)

	return &Client{firecracker.NewClient(socketPath, logrus.NewEntry(logs), true)}
}

// GetResource returns the machine configuration. eg. CPU, Memory, etc.
func (c *Client) GetResource() (*models.MachineConfiguration, error) {
	cfg, err := c.Client.GetMachineConfiguration()
	if err != nil {
		return nil, fmt.Errorf("failed to get machine resource configuration: %w", err)
	}
	return cfg.GetPayload(), nil
}

// GetInstance returns the machine instance.
func (c *Client) GetInstance(ctx context.Context) (*models.InstanceInfo, error) {
	inst, err := c.Client.GetInstanceInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get machine instance instance: %w", err)
	}
	return inst.GetPayload(), nil
}

// GetVmConfig returns the machine configuration.
func (c *Client) GetVmConfig(ctx context.Context) (*models.FullVMConfiguration, error) {
	cfg, err := c.Client.GetExportVMConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get machine vm configuration: %w", err)
	}
	return cfg.GetPayload(), nil
}
