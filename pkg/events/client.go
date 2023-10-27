package events

import "context"

// Client is the interface for the pub/sub client
type Client interface {
	// Publish publishes a message to the broker
	Publish(context.Context, string, []byte) error
	// Subscribe subscribes to a topic
	Subscribe(context.Context, string) (chan []byte, error)
	// Unsubscribe unsubscribes from a topic
	Unsubscribe(context.Context, string) error
	// Connect connects to the broker
	Connect(context.Context) error
	// Disconnect disconnects from the broker
	Disconnect() error
}
