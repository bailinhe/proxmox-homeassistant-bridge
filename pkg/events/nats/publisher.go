package nats

import (
	"context"
	"strings"

	"go.uber.org/zap"
)

// Publish publishes a message to the broker
func (c *Client) Publish(_ context.Context, topic string, msg []byte) error {
	topic = strings.ReplaceAll(topic, "/", ".")
	c.logger.Debug("publishing message", zap.String("topic", topic))

	return c.nc.Publish(topic, msg)
}
