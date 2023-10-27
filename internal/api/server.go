package api

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/bighotel/proxmox-ha-bridge/pkg/events"
	"gitlab.com/bighotel/proxmox-ha-bridge/pkg/proxmox"
	"gitlab.com/bighotel/proxmox-ha-bridge/pkg/vmcontroller"
	"go.uber.org/zap"
)

// Server is the main server
type Server struct {
	pc         *proxmox.Client
	ec         events.Client
	monitorVMs map[string]uint8

	logger *zap.Logger

	probeInterval               *time.Duration
	availabilityPublishInterval *time.Duration
}

// Opt is a server option
type Opt func(*Server)

// NewServer creates a new server
func NewServer(monitorVMIDs []int, opts ...Opt) *Server {
	s := &Server{
		logger:     zap.NewNop(),
		monitorVMs: make(map[string]uint8),
	}

	for _, opt := range opts {
		opt(s)
	}

	s.logger = s.logger.With(zap.String("component", "api-server"))

	for _, vmid := range monitorVMIDs {
		s.monitorVMs[fmt.Sprint(vmid)] = 0
	}

	return s
}

// WithProbeInterval sets the probe interval
func WithProbeInterval(d time.Duration) Opt {
	return func(s *Server) {
		s.probeInterval = &d
	}
}

// WithAvailabilityPublishInterval sets the availability publish interval
func WithAvailabilityPublishInterval(d time.Duration) Opt {
	return func(s *Server) {
		s.availabilityPublishInterval = &d
	}
}

// WithProxmoxClient sets the proxmox client
func WithProxmoxClient(p *proxmox.Client) Opt {
	return func(s *Server) {
		s.pc = p
	}
}

// WithLogger sets the logger
func WithLogger(logger *zap.Logger) Opt {
	return func(s *Server) {
		s.logger = logger
	}
}

// WithEventsClient sets the events client
func WithEventsClient(ec events.Client) Opt {
	return func(s *Server) {
		s.ec = ec
	}
}

func (s *Server) scanVMs() ([]vmcontroller.VMController, error) {
	vms, err := s.pc.VMs()
	if err != nil {
		return nil, err
	}

	filteredVMs := make([]vmcontroller.VMController, 0)

	for _, vm := range vms {
		if _, ok := s.monitorVMs[vm.ID()]; ok {
			filteredVMs = append(filteredVMs, vm)
		}
	}

	return filteredVMs, nil
}

// Start starts the server
func (s *Server) Start() error {
	vms, err := s.scanVMs()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := s.ec.Connect(ctx); err != nil {
		return err
	}

	defer func() {
		err := s.ec.Disconnect()
		if err != nil {
			s.logger.Error("error disconnecting from events client", zap.Error(err))
		}
	}()

	vmcs := make([]*vmcontroller.VirtualMachineMQTTDevice, len(vms))

	for i, vm := range vms {
		s.logger.Debug("monitoring vm", zap.String("vmid", vm.ID()))

		opts := []vmcontroller.MQTTDeviceOpt{
			vmcontroller.WithLogger(s.logger),
		}

		if s.probeInterval != nil {
			opts = append(opts, vmcontroller.WithProbeInterval(*s.probeInterval))
		}

		if s.availabilityPublishInterval != nil {
			opts = append(opts, vmcontroller.WithAvailabilityPublishInterval(*s.availabilityPublishInterval))
		}

		vmcs[i] = vmcontroller.NewVirtualMachineMQTTDevice(vm, s.ec, opts...)
		go vmcs[i].Run(context.Background())
	}

	// listen to SIGINT, SIGTERM and stop gracefully
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh

	s.logger.Info("shutting down server gracefully")

	cancel()

	for _, vmc := range vmcs {
		vmc.Stop()
	}

	// wait one second and exit
	time.Sleep(1 * time.Second)

	return nil
}
