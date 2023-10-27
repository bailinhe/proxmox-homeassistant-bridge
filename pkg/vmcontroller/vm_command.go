package vmcontroller

import (
	"context"
	"encoding/json"
	"time"

	"gitlab.com/bighotel/proxmox-ha-bridge/pkg/events"
	"gitlab.com/bighotel/proxmox-ha-bridge/pkg/homeassistant"
	"go.uber.org/zap"
)

// VMCommandStateMessage is the state message for a virtual machine command
type VMCommandStateMessage struct {
	State string `json:"state"`
}

// Marshal marshals the state message
func (m *VMCommandStateMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// Unmarshal unmarshals the state message
func (m *VMCommandStateMessage) Unmarshal(msg []byte) error {
	return json.Unmarshal(msg, m)
}

// VMCommandStateMessage implements MQTTMessage
var _ homeassistant.MQTTMessage = &VMCommandStateMessage{}

// VMCommandMessage is the message for a virtual machine command
type VMCommandMessage struct {
	Command VMCommandOption `json:"command"`
}

// Marshal marshals the command message
func (m *VMCommandMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// Unmarshal unmarshals the command message
func (m *VMCommandMessage) Unmarshal(msg []byte) error {
	return json.Unmarshal(msg, m)
}

// VMCommandMessage implements MQTTMessage
var _ homeassistant.MQTTMessage = &VMCommandMessage{}

// VMCommandOption is the enum type for a virtual machine command option
type VMCommandOption string

const (
	// VMCommandOptionStart is the start option for a virtual machine command
	VMCommandOptionStart VMCommandOption = "start"
	// VMCommandOptionStop is the stop option for a virtual machine command
	VMCommandOptionStop VMCommandOption = "stop"
	// VMCommandOptionReboot is the reboot option for a virtual machine command
	VMCommandOptionReboot VMCommandOption = "reboot"
	// VMCommandOptionReset is the reset option for a virtual machine command
	VMCommandOptionReset VMCommandOption = "reset"
	// VMCommandOptionShutdown is the shutdown option for a virtual machine command
	VMCommandOptionShutdown VMCommandOption = "shutdown"
)

// VMCommandSelect is a home assistant select entity that sends commands to a
// virtual machine
type VMCommandSelect struct {
	homeassistant.MQTTSelect

	vm                          VMController
	availabilityPublishInterval time.Duration
	device                      *homeassistant.MQTTDevice
	logger                      *zap.Logger
	stopchan                    chan struct{}
}

// VMCommandSelectOpt is a functional option for a VMCommandSelect
type VMCommandSelectOpt func(*VMCommandSelect)

// VMCommandSelectWithAvailabilityPublishInterval sets the availability publish
// interval for the VMCommandSelect
func VMCommandSelectWithAvailabilityPublishInterval(d time.Duration) VMCommandSelectOpt {
	return func(s *VMCommandSelect) {
		s.availabilityPublishInterval = d
	}
}

// VMCommandSelectWithLogger sets the logger for the VMCommandSelect
func VMCommandSelectWithLogger(logger *zap.Logger) VMCommandSelectOpt {
	return func(s *VMCommandSelect) {
		s.logger = logger
	}
}

// NewVMCommandSelect creates a new VMCommandSelect
func NewVMCommandSelect(
	vm VMController,
	ec events.Client,
	device *homeassistant.MQTTDevice,
	opts ...VMCommandSelectOpt,
) *VMCommandSelect {
	s := &VMCommandSelect{
		availabilityPublishInterval: DefaultAvailabilityPublishInterval,
		logger:                      zap.NewNop(),
		vm:                          vm,
		device:                      device,

		stopchan: make(chan struct{}),
	}

	cfgtopic, cfg := mkSelectConfig("Command", "vm-command", device)
	cfg.Icon = "mdi:light-switch-off"

	for _, opt := range opts {
		opt(s)
	}

	s.logger = s.logger.With(zap.String("component", "vm-command-select"))

	s.MQTTSelect = *homeassistant.NewMQTTSelect(
		cfg, cfgtopic, ec,
		homeassistant.MQTTSelectWithLogger(s.logger),
	)

	return s
}

// Run runs the VMCommandSelect
func (s *VMCommandSelect) Run(ctx context.Context) {
	s.logger.Info("starting vm command select", zap.String("vm", s.vm.Name()))
	defer s.logger.Info("vm command select stopped", zap.String("vm", s.vm.Name()))

	if err := s.PublishConfig(ctx); err != nil {
		s.logger.Error("failed to publish config", zap.Error(err), zap.String("vm", s.vm.Name()))
	}

	if err := s.PublishAvailability(ctx, &VMStateSensorAvailabilityMessage{
		Status: VMStateAvailabilityOnline,
	}); err != nil {
		s.logger.Error("failed to publish availability", zap.Error(err), zap.String("vm", s.vm.Name()))
	}

	availabilityticker := time.NewTicker(s.availabilityPublishInterval)
	defer availabilityticker.Stop()

	msgchan, err := s.SubscribeCommand(ctx)
	if err != nil {
		s.logger.Error("failed to subscribe to command topic", zap.Error(err), zap.String("vm", s.vm.Name()))
		return
	}

	defer func() {
		err := s.UnsubscribeCommand(ctx)
		if err != nil {
			s.logger.Error(
				"failed to unsubscribe from command topic",
				zap.Error(err), zap.String("vm", s.vm.Name()),
			)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("context done", zap.String("vm", s.vm.Name()))
			return
		case <-s.stopchan:
			return
		case <-availabilityticker.C:
			if err := s.PublishAvailability(ctx, &VMStateSensorAvailabilityMessage{
				Status: VMStateAvailabilityOnline,
			}); err != nil {
				s.logger.Error("failed to publish availability", zap.Error(err), zap.String("vm", s.vm.Name()))
			}
		case msg := <-msgchan:
			cmd := &VMCommandMessage{}

			err := cmd.Unmarshal(msg)
			if err != nil {
				s.logger.Error("failed to unmarshal command message", zap.Error(err), zap.String("vm", s.vm.Name()))
				continue
			}

			s.logger.Info("received command", zap.String("command", string(cmd.Command)), zap.String("vm", s.vm.Name()))

			func() {
				if err = s.handleCommand(ctx, cmd.Command); err != nil {
					s.logger.Error("handle command error", zap.Error(err), zap.String("vm", s.vm.Name()))
				}
			}()
		}
	}
}

// Stop stops the VMCommandSelect
func (s *VMCommandSelect) Stop() {
	s.logger.Info("stopping vm command select", zap.String("vm", s.vm.Name()))

	s.stopchan <- struct{}{}

	if err := s.PublishAvailability(context.Background(), &VMStateSensorAvailabilityMessage{
		Status: VMStateAvailabilityOffline,
	}); err != nil {
		s.logger.Error("failed to publish availability", zap.Error(err), zap.String("vm", s.vm.Name()))
	}
}

// handleCommand handles a command message
func (s *VMCommandSelect) handleCommand(ctx context.Context, cmd VMCommandOption) error {
	var err error

	switch cmd {
	case VMCommandOptionStart:
		err = s.vm.Start(ctx)
	case VMCommandOptionStop:
		err = s.vm.Stop(ctx)
	case VMCommandOptionReboot:
		err = s.vm.Restart(ctx)
	case VMCommandOptionReset:
		err = s.vm.Reset(ctx)
	case VMCommandOptionShutdown:
		err = s.vm.Shutdown(ctx)
	}

	return err
}

// VMCommandSelect implements homeassistant.MQTTDevice
var _ Select = (*VMCommandSelect)(nil)
