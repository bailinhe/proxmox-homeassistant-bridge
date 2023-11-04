package vmcontroller

import (
	"context"

	"go.uber.org/zap"
)

// VMController is a virtual machine controller
type VMController interface {
	// Start starts the vm
	Start(context.Context) error
	// Stop stops the vm
	Stop(context.Context) error
	// Restart restarts the vm
	Restart(context.Context) error
	// Reset hard resets the vm
	Reset(context.Context) error
	// Shutdown gracefully shuts down the vm
	Shutdown(context.Context) error
	// Status returns the status of the vm
	Status() (string, error)
	// ID returns the id of the vm
	ID() string
	// Name returns the name of the vm
	Name() string
	// SetLogger sets the logger
	SetLogger(logger *zap.Logger)
}
