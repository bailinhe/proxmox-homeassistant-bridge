package mqtt

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gitlab.com/bighotel/proxmox-ha-bridge/pkg/events"
	"go.uber.org/zap"
)

// Client is a wrapper around the MQTT client
type Client struct {
	client mqtt.Client
	logger *zap.Logger
}

// Opt is a functional option for the Client
type Opt func(c *Client)

// NewClient creates a new Client
func NewClient(opts ...Opt) *Client {
	c := &Client{}

	for _, opt := range opts {
		opt(c)
	}

	c.logger = c.logger.With(zap.String("component", "mqtt-client"))

	return c
}

// WithLogger sets the logger for the Client
func WithLogger(l *zap.Logger) Opt {
	return func(c *Client) {
		c.logger = l
	}
}

// WithClient sets the MQTT client for the Client
func WithClient(client mqtt.Client) Opt {
	return func(c *Client) {
		c.client = client
	}
}

// Connect connects to the MQTT broker
func (c *Client) Connect(ctx context.Context) error {
	if c.client == nil {
		return ErrNoClient
	}

	if c.client.IsConnected() {
		return nil
	}

	token := c.client.Connect()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-token.Done():
		if token.Error() != nil {
			return token.Error()
		}
	}

	return nil
}

// Disconnect disconnects from the MQTT broker
func (c *Client) Disconnect() error {
	if c.client == nil {
		return ErrNoClient
	}

	if !c.client.IsConnected() {
		return nil
	}

	disconnectWaitMS := 250
	c.client.Disconnect(uint(disconnectWaitMS))

	return nil
}

// Client implements the Client interface
var _ events.Client = (*Client)(nil)
