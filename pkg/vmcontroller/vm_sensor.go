package vmcontroller

import (
	"context"
	"encoding/json"
	"time"

	"gitlab.com/bighotel/proxmox-ha-bridge/pkg/events"
	"gitlab.com/bighotel/proxmox-ha-bridge/pkg/homeassistant"
	"go.uber.org/zap"
)

// VMStateSensorStateMessage is the state message for a virtual machine state sensor
type VMStateSensorStateMessage struct {
	State string `json:"state"`
}

// Marshal marshals the state message
func (m *VMStateSensorStateMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// Unmarshal unmarshals the state message
func (m *VMStateSensorStateMessage) Unmarshal(msg []byte) error {
	return json.Unmarshal(msg, m)
}

// VMStateAvailability is the availability enum type for a virtual machine state sensor
type VMStateAvailability string

const (
	// VMStateAvailabilityOnline is the online availability for a virtual machine state sensor
	VMStateAvailabilityOnline VMStateAvailability = "online"
	// VMStateAvailabilityOffline is the offline availability for a virtual machine state sensor
	VMStateAvailabilityOffline VMStateAvailability = "offline"
)

// VMStateSensorAvailabilityMessage is the availability message for a virtual machine state sensor
type VMStateSensorAvailabilityMessage struct {
	Status VMStateAvailability `json:"status"`
}

// Marshal marshals the availability message
func (m *VMStateSensorAvailabilityMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// Unmarshal unmarshals the availability message
func (m *VMStateSensorAvailabilityMessage) Unmarshal(msg []byte) error {
	return json.Unmarshal(msg, m)
}

// VMStateSensorAvailabilityMessage implements MQTTMessage
var _ homeassistant.MQTTMessage = &VMStateSensorAvailabilityMessage{}

// VMStateSensorStateMessage implements MQTTMessage
var _ homeassistant.MQTTMessage = &VMStateSensorStateMessage{}

// VMStateSensor is a home assistant sensor that reports the state of a virtual machine
type VMStateSensor struct {
	homeassistant.MQTTSensor

	vm VMController

	probeInterval               time.Duration
	availabilityPublishInterval time.Duration
	device                      *homeassistant.MQTTDevice
	logger                      *zap.Logger

	stopchan chan struct{}
}

// VMStateSensorOpt is a functional option for the VMStateSensor
type VMStateSensorOpt func(*VMStateSensor)

// VMStateSensorWithProbeInterval sets the probe interval for the VMStateSensor
func VMStateSensorWithProbeInterval(d time.Duration) VMStateSensorOpt {
	return func(s *VMStateSensor) {
		s.probeInterval = d
	}
}

// VMStateSensorWithAvailabilityPublishInterval sets the availability publish interval for the VMStateSensor
func VMStateSensorWithAvailabilityPublishInterval(d time.Duration) VMStateSensorOpt {
	return func(s *VMStateSensor) {
		s.availabilityPublishInterval = d
	}
}

// VMStateSensorWithLogger sets the logger for the VMStateSensor
func VMStateSensorWithLogger(logger *zap.Logger) VMStateSensorOpt {
	return func(s *VMStateSensor) {
		s.logger = logger
	}
}

// NewVMStateSensor creates a new virtual machine state sensor
func NewVMStateSensor(
	vm VMController,
	ec events.Client,
	device *homeassistant.MQTTDevice,
	opts ...VMStateSensorOpt,
) *VMStateSensor {
	s := &VMStateSensor{
		vm:                          vm,
		device:                      device,
		probeInterval:               DefaultProbeInterval,
		availabilityPublishInterval: DefaultAvailabilityPublishInterval,

		stopchan: make(chan struct{}),
	}

	cfgtopic, cfg := mkSensorConfig("VM Status", "vm-status", s.device)
	cfg.Icon = "mdi:server"

	for _, opt := range opts {
		opt(s)
	}

	s.logger = s.logger.With(zap.String("component", "vm-state-sensor"))

	s.MQTTSensor = *homeassistant.NewMQTTSensor(
		cfg, cfgtopic, ec,
		homeassistant.MQTTSensorWithLogger(s.logger),
	)

	return s
}

// UpdateAndPublishState updates and publishes the state of the sensor
func (s *VMStateSensor) updateAndPublishState(ctx context.Context) error {
	state, err := s.vm.Status()
	if err != nil {
		return err
	}

	msg := &VMStateSensorStateMessage{
		State: state,
	}

	return s.PublishState(ctx, msg)
}

// Run runs the VMStateSensor
func (s *VMStateSensor) Run(ctx context.Context) {
	s.logger.Info("starting vm state sensor", zap.String("vm", s.vm.Name()))
	defer s.logger.Info("vm state sensor stopped", zap.String("vm", s.vm.Name()))

	if err := s.PublishConfig(ctx); err != nil {
		s.logger.Error("failed to publish config", zap.Error(err))
	}

	if err := s.PublishAvailability(ctx, &VMStateSensorAvailabilityMessage{
		Status: VMStateAvailabilityOnline,
	}); err != nil {
		s.logger.Error("failed to publish availability", zap.Error(err))
	}

	if err := s.updateAndPublishState(ctx); err != nil {
		s.logger.Error("failed to update and publish state", zap.Error(err))
	}

	probeticker := time.NewTicker(s.probeInterval)
	defer probeticker.Stop()

	availabilityticker := time.NewTicker(s.availabilityPublishInterval)
	defer availabilityticker.Stop()

	for {
		select {
		case <-s.stopchan:
			return
		case <-ctx.Done():
			s.logger.Info("context done, stopping vm state sensor", zap.String("vm", s.vm.Name()))
			return
		case <-probeticker.C:
			if err := s.updateAndPublishState(ctx); err != nil {
				s.logger.Error("failed to update and publish state", zap.Error(err))
			}
		case <-availabilityticker.C:
			if err := s.PublishAvailability(ctx, &VMStateSensorAvailabilityMessage{
				Status: VMStateAvailabilityOnline,
			}); err != nil {
				s.logger.Error("failed to publish availability", zap.Error(err))
			}

			if err := s.PublishConfig(ctx); err != nil {
				s.logger.Error("failed to publish config", zap.Error(err))
			}
		}
	}
}

// Stop stops the VMStateSensor
func (s *VMStateSensor) Stop() {
	s.logger.Info("stopping vm state sensor", zap.String("vm", s.vm.Name()))

	s.stopchan <- struct{}{}

	if err := s.PublishAvailability(context.Background(), &VMStateSensorAvailabilityMessage{
		Status: VMStateAvailabilityOffline,
	}); err != nil {
		s.logger.Error("failed to publish availability", zap.Error(err))
	}
}

// VMStateSensor implements Sensor
var _ Sensor = (*VMStateSensor)(nil)
