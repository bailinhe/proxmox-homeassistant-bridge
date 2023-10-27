package nats

import (
	"context"
	"strings"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// Subscribe subscribes to a topic
func (c *Client) Subscribe(_ context.Context, topic string) (chan []byte, error) {
	if c.nc == nil {
		return nil, ErrNoClient
	}

	topic = strings.ReplaceAll(topic, "/", ".")

	msgChanCapcity := 64
	msgchan := make(chan []byte, msgChanCapcity)

	c.logger.Debug("subscribing to topic", zap.String("topic", topic))

	sub, err := c.nc.Subscribe(topic, func(m *nats.Msg) {
		msgchan <- m.Data
	})
	if err != nil {
		return nil, err
	}

	c.subs[topic] = sub

	return msgchan, nil
}

// Unsubscribe unsubscribes from a topic
func (c *Client) Unsubscribe(_ context.Context, topic string) error {
	if c.nc == nil {
		return ErrNoClient
	}

	topic = strings.ReplaceAll(topic, "/", ".")

	sub, ok := c.subs[topic]
	if !ok {
		return ErrNoSubscription
	}

	c.logger.Debug("unsubscribing from topic", zap.String("topic", topic))

	err := sub.Unsubscribe()
	if err != nil {
		return err
	}

	delete(c.subs, topic)

	return nil
}
