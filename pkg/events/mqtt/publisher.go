package mqtt

import "context"

// Publish publishes a message to the broker
func (c *Client) Publish(ctx context.Context, topic string, message []byte) error {
	if c.client == nil {
		return ErrNoClient
	}

	t := c.client.Publish(topic, 0, true, message)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.Done():
		if t.Error() != nil {
			return t.Error()
		}
	}

	c.logger.Sugar().Debugf("published message: %s", message)

	return nil
}
