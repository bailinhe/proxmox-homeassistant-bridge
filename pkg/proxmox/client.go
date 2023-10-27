package proxmox

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/luthermonson/go-proxmox"
	"go.uber.org/zap"
)

// Client is a wrapper around the Proxmox client
type Client struct {
	serverURL string

	client     *proxmox.Client
	clientOpts []proxmox.Option
	logger     *zap.Logger

	nodes proxmox.NodeStatuses
}

// Opt is a functional option for the Client
type Opt func(c *Client)

// NewClient creates a new Client
func NewClient(serverURL string, opts ...Opt) *Client {
	c := &Client{
		serverURL: serverURL,
		logger:    &zap.Logger{},

		clientOpts: []proxmox.Option{},
	}

	for _, opt := range opts {
		opt(c)
	}

	c.logger = c.logger.With(zap.String("component", "proxmox"))
	c.client = proxmox.NewClient(
		fmt.Sprintf("%s/api2/json", c.serverURL),
		c.clientOpts...,
	)

	return c
}

// WithAPIToken sets the API token for the Client
func WithAPIToken(tokenID, secret string) Opt {
	return func(c *Client) {
		c.clientOpts = append(c.clientOpts, proxmox.WithAPIToken(tokenID, secret))
	}
}

// WithLogger sets the logger for the Client
func WithLogger(l *zap.Logger) Opt {
	return func(c *Client) {
		c.logger = l
	}
}

// WithInsecure sets the insecure flag for the Client
func WithInsecure() Opt {
	return func(c *Client) {
		c.clientOpts = append(c.clientOpts, proxmox.WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}))
	}
}

// Version returns the version of the Proxmox server
func (c *Client) Version() (*proxmox.Version, error) {
	return c.client.Version()
}

// ProxmoxClient returns the Proxmox client
func (c *Client) ProxmoxClient() *proxmox.Client {
	return c.client
}

// Start starts the Client
func (c *Client) Start() {
	NodeStatuses, err := c.client.Nodes()
	if err != nil {
		c.logger.Error(err.Error())
		return
	}

	vms := proxmox.VirtualMachines{}

	for _, st := range NodeStatuses {
		node, err := c.client.Node(st.Node)
		if err != nil {
			c.logger.Error(err.Error())
			continue
		}

		vm, err := node.VirtualMachines()
		if err != nil {
			c.logger.Error(err.Error())
			continue
		}

		vms = append(vms, vm...)
	}

	for _, vm := range vms {
		j, _ := json.MarshalIndent(vm, "", "  ")
		c.logger.Sugar().Debug(string(j))
	}
}
