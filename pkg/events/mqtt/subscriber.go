package mqtt

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

// Subscribe subscribes to a topic
func (c *Client) Subscribe(ctx context.Context, topic string) (chan []byte, error) {
	if c.client == nil {
		return nil, ErrNoClient
	}

	msgchan := make(chan []byte)

	t := c.client.Subscribe(topic, 1, func(_ mqtt.Client, m mqtt.Message) { // nolint: gomnd
		c.logger.Debug("received message", zap.String("topic", m.Topic()), zap.String("payload", string(m.Payload())))
		msgchan <- m.Payload()
		c.logger.Debug("sent message to channel")
	})

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-t.Done():
		if t.Error() != nil {
			return nil, t.Error()
		}
	}

	c.logger.Sugar().Infof("subscribed to topic: %s", topic)

	return msgchan, nil
}

// Unsubscribe unsubscribes from a topic
func (c *Client) Unsubscribe(ctx context.Context, topic string) error {
	if c.client == nil {
		return ErrNoClient
	}

	t := c.client.Unsubscribe(topic)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.Done():
		if t.Error() != nil {
			return t.Error()
		}
	}

	c.logger.Sugar().Infof("unsubscribed from topic: %s", topic)

	return nil
}
