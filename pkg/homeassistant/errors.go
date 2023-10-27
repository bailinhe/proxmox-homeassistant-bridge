package homeassistant

import "errors"

var (
	// ErrNoClient is the error returned when there is no MQTT client
	ErrNoClient = errors.New("no MQTT client")
	// ErrNoAvailabilityTopic is the error returned when there is no availability topic
	ErrNoAvailabilityTopic = errors.New("no availability topic")
	// ErrNoStateTopic is the error returned when there is no state topic
	ErrNoStateTopic = errors.New("no state topic")
	// ErrNoCommandTopic is the error returned when there is no command topic
	ErrNoCommandTopic = errors.New("no command topic")
)
