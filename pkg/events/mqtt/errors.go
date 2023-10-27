package mqtt

import "errors"

// ErrNoClient is returned when the client is nil
var ErrNoClient = errors.New("no client provided")
