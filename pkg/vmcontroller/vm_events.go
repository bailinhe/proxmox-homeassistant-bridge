package vmcontroller

import (
	"context"
	"encoding/json"
	"time"

	"gitlab.com/bighotel/proxmox-ha-bridge/pkg/events"
	"gitlab.com/bighotel/proxmox-ha-bridge/pkg/homeassistant"
	"go.uber.org/zap"
)

// VMEventSensorEventMessage is the event message for a virtual machine event
// sensor
type VMEventSensorEventMessage struct {
	Event string `json:"event"`
}

// Marshal marshals the event message
func (m *VMEventSensorEventMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// Unmarshal unmarshals the event message
func (m *VMEventSensorEventMessage) Unmarshal(msg []byte) error {
	return json.Unmarshal(msg, m)
}

// VMEventSensorEventMessage implements MQTTMessage
var _ homeassistant.MQTTMessage = &VMEventSensorEventMessage{}

// VMEventSensor is a home assistant sensor that reports the events of a
// virtual machine
type VMEventSensor struct {
	*homeassistant.MQTTSensor

	availabilityPublishInterval time.Duration
	device                      *homeassistant.MQTTDevice
	logger                      *zap.Logger
	stopchan                    chan struct{}
}

// VMEventSensorOpt is a functional option for configuring a VMEventSensor
type VMEventSensorOpt func(*VMEventSensor)

// VMEventSensorWithAvailabilityPublishInterval sets the availability publish
// interval for a VMEventSensor
func VMEventSensorWithAvailabilityPublishInterval(d time.Duration) VMEventSensorOpt {
	return func(s *VMEventSensor) {
		s.availabilityPublishInterval = d
	}
}

// VMEventSensorWithLogger sets the logger for a VMEventSensor
func VMEventSensorWithLogger(l *zap.Logger) VMEventSensorOpt {
	return func(s *VMEventSensor) {
		s.logger = l
	}
}

// NewVMEventSensor creates a new VMEventSensor
func NewVMEventSensor(
	ec events.Client, device *homeassistant.MQTTDevice, opts ...VMEventSensorOpt,
) *VMEventSensor {
	s := &VMEventSensor{
		device:                      device,
		availabilityPublishInterval: DefaultAvailabilityPublishInterval,

		stopchan: make(chan struct{}),
	}

	cfgtopic, cfg := mkSensorConfig("VM Events", "vm-events", s.device)

	cfg.Icon = "mdi:message-bulleted"
	cfg.ValueTemplate = "{{ value_json.event }}"

	for _, opt := range opts {
		opt(s)
	}

	s.logger = s.logger.With(
		zap.String("component", "vm-events"),
		zap.String("name", s.device.Name),
	)

	s.MQTTSensor = homeassistant.NewMQTTSensor(
		cfg, cfgtopic, ec,
		homeassistant.MQTTSensorWithLogger(s.logger),
	)

	return s
}

// Run starts the VMEventSensor
func (s *VMEventSensor) Run(ctx context.Context) {
	s.logger.Info("starting vm event sensor")
	defer s.logger.Info("stopped vm event sensor")

	if err := s.PublishConfig(ctx); err != nil {
		s.logger.Error("failed to publish config", zap.Error(err))
	}

	if err := s.PublishAvailability(ctx, &VMStateSensorAvailabilityMessage{
		Status: VMStateAvailabilityOnline,
	}); err != nil {
		s.logger.Error("failed to publish availability", zap.Error(err))
	}

	availabilityticker := time.NewTicker(s.availabilityPublishInterval)
	defer availabilityticker.Stop()

	for {
		select {
		case <-s.stopchan:
			return
		case <-ctx.Done():
			s.logger.Info("context done, stopping vm event sensor")
			return
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

// Stop stops the VMEventSensor
func (s *VMEventSensor) Stop() {
	s.logger.Info("stopping vm event sensor")
	s.stopchan <- struct{}{}

	if err := s.PublishAvailability(context.Background(), &VMStateSensorAvailabilityMessage{
		Status: VMStateAvailabilityOffline,
	}); err != nil {
		s.logger.Error("failed to publish availability", zap.Error(err))
	}
}

// EmitEvent emits an event
func (s *VMEventSensor) EmitEvent(ctx context.Context, payload string) error {
	event := &VMEventSensorEventMessage{
		Event: payload,
	}

	s.logger.Debug("emitting event", zap.String("payload", payload))

	return s.PublishState(ctx, event)
}

// VMEventSensor implements Sensor
var _ Sensor = (*VMEventSensor)(nil)
