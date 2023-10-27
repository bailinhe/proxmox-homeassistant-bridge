package nats

import (
	"context"

	"github.com/nats-io/nats.go"
	"gitlab.com/bighotel/proxmox-ha-bridge/pkg/events"
	"go.uber.org/zap"
)

// Client is a wrapper around the NATS client
type Client struct {
	nc       *nats.Conn
	natsurl  string
	natsopts []nats.Option

	subs map[string]*nats.Subscription

	logger *zap.Logger
}

// Opt is a functional option for the Client
type Opt func(c *Client)

// NewClient creates a new Client
func NewClient(opts ...Opt) *Client {
	c := &Client{
		subs:   make(map[string]*nats.Subscription),
		logger: zap.NewNop(),
	}

	for _, opt := range opts {
		opt(c)
	}

	c.logger = c.logger.With(zap.String("component", "nats-client"))

	return c
}

// WithLogger sets the logger for the Client
func WithLogger(l *zap.Logger) Opt {
	return func(c *Client) {
		c.logger = l
	}
}

// WithNATSURL sets the NATS URL for the Client
func WithNATSURL(url string) Opt {
	return func(c *Client) {
		c.natsurl = url
	}
}

// WithNATSOpts sets the NATS options for the Client
func WithNATSOpts(opts ...nats.Option) Opt {
	return func(c *Client) {
		c.natsopts = opts
	}
}

// Connect connects to the NATS broker
func (c *Client) Connect(_ context.Context) error {
	nc, err := nats.Connect(c.natsurl, c.natsopts...)
	if err != nil {
		return err
	}

	c.nc = nc

	if c.nc.IsConnected() {
		return nil
	}

	return nil
}

// Disconnect disconnects from the NATS broker
func (c *Client) Disconnect() error {
	if c.nc == nil {
		return nil
	}

	c.nc.Close()

	return nil
}

// Client implements the events.Client interface
var _ events.Client = (*Client)(nil)
