package nats

import "errors"

// ErrNoClient is returned when the client is nil
var ErrNoClient = errors.New("no client provided")

// ErrNoSubscription is returned when the subscription is nil
var ErrNoSubscription = errors.New("not subscribed to topic")
