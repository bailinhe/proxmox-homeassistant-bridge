package homeassistant

import (
	"context"
	"encoding/json"

	"gitlab.com/bighotel/proxmox-ha-bridge/pkg/events"
	"go.uber.org/zap"
)

// MQTTSelectConfiguration is the MQTT Select type.
// The mqtt Select platform allows you to integrate devices that might expose
// configuration options through MQTT into Home Assistant as a Select. Every
// time a message under the topic in the configuration is received, the select
// entity will be updated in Home Assistant and vice-versa, keeping the device
// and Home Assistant in sync.
type MQTTSelectConfiguration struct {
	// availability list (optional) A list of MQTT topics subscribed to receive
	// availability (online/offline) updates. Must not be used together with
	// availability_topic.
	Availability []MQTTAvailability `json:"availability,omitempty"`
	// availability_topic string (optional) The MQTT topic subscribed to receive
	// availability (online/offline) updates. Must not be used together with
	// availability.
	AvailabilityTopic string `json:"availability_topic,omitempty"`
	// availability_mode string (optional, default: latest) When availability is
	// configured, this controls the conditions needed to set the entity to
	// available. Valid entries are all, any, and latest. If set to all,
	// payload_available must be received on all configured availability topics
	// before the entity is marked as online. If set to any, payload_available
	// must be received on at least one configured availability topic before the
	// entity is marked as online. If set to latest, the last payload_available
	// or payload_not_available received on any configured availability topic
	// controls the availability.
	AvailabilityMode MQTTAvailabilityMode `json:"availability_mode,omitempty"`
	// availability_template template (optional) Defines a template to extract
	// device’s availability from the availability_topic. To determine the
	// devices’s availability result of this template will be compared to
	// payload_available and payload_not_available.
	AvailabilityTemplate string `json:"availability_template,omitempty"`
	// command_template template (optional) Defines a template to generate the
	// payload to send to command_topic.
	CommandTemplate string `json:"command_template,omitempty"`
	// command_topic string REQUIRED The MQTT topic to publish commands to change
	// the selected option.
	CommandTopic string `json:"command_topic"`
	// device map (optional) Information about the device this Select is a part
	// of to tie it into the device registry. Only works when unique_id is set.
	// At least one of identifiers or connections must be present to identify
	// the device.
	Device *MQTTDevice `json:"device,omitempty"`
	// enabled_by_default boolean (optional, default: true) Flag which defines
	// if the entity should be enabled when first added.
	EnabledByDefault bool `json:"enabled_by_default,omitempty"`
	// encoding string (optional, default: utf-8) The encoding of the payloads
	// received and published messages. Set to "" to disable decoding of incoming
	// payload.
	Encoding string `json:"encoding,omitempty"`
	// entity_category string (optional, default: None) The category of the entity.
	EntityCategory string `json:"entity_category,omitempty"`
	// icon icon (optional) Icon for the entity.
	Icon string `json:"icon,omitempty"`
	// json_attributes_template template (optional) Defines a template to extract
	// the JSON dictionary from messages received on the json_attributes_topic.
	JSONAttributesTemplate string `json:"json_attributes_template,omitempty"`
	// json_attributes_topic string (optional) The MQTT topic subscribed to
	// receive a JSON dictionary payload and then set as entity attributes.
	// Implies force_update of the current select state when a message is
	// received on this topic.
	JSONAttributesTopic string `json:"json_attributes_topic,omitempty"`
	// name string (optional) The name of the Select. Can be set to null if only
	// the device name is relevant.
	Name string `json:"name,omitempty"`
	// object_id string (optional) Used instead of name for automatic generation
	// of entity_id
	ObjectID string `json:"object_id,omitempty"`
	// optimistic boolean (optional) Flag that defines if the select works in
	// optimistic mode. Default: true if no state_topic defined, else false.
	Optimistic bool `json:"optimistic,omitempty"`
	// options list REQUIRED List of options that can be selected. An empty list
	// or a list with a single item is allowed.
	Options []string `json:"options"`
	// qos integer (optional, default: 0) The maximum QoS level to be used when
	// receiving and publishing messages.
	QoS int `json:"qos,omitempty"`
	// retain boolean (optional, default: false) If the published message should
	// have the retain flag on or not.
	Retain bool `json:"retain,omitempty"`
	// state_topic string (optional) The MQTT topic subscribed to receive update
	// of the selected option.
	StateTopic string `json:"state_topic,omitempty"`
	// unique_id string (optional) An ID that uniquely identifies this Select.
	// If two Selects have the same unique ID Home Assistant will raise an
	// exception.
	UniqueID string `json:"unique_id,omitempty"`
	// value_template template (optional) Defines a template to extract the value.
	ValueTemplate string `json:"value_template,omitempty"`
}

// MQTTSelect is a MQTT Select entity.
type MQTTSelect struct {
	Configuration *MQTTSelectConfiguration
	ConfigTopic   string
	ec            events.Client
	logger        *zap.Logger
}

// MQTTSelectOpt is a functional option for the MQTTSelect.
type MQTTSelectOpt func(*MQTTSelect)

// MQTTSelectWithLogger sets the logger for the MQTTSelect
func MQTTSelectWithLogger(logger *zap.Logger) MQTTSelectOpt {
	return func(s *MQTTSelect) {
		s.logger = logger
	}
}

// NewMQTTSelect creates a new MQTTSelect.
func NewMQTTSelect(cfg *MQTTSelectConfiguration, cfgtopic string, ec events.Client, opts ...MQTTSelectOpt) *MQTTSelect {
	s := &MQTTSelect{
		Configuration: cfg,
		ConfigTopic:   cfgtopic,
		ec:            ec,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.logger = s.logger.With(zap.String("component", "mqtt-select"))

	return s
}

// PublishAvailability publishes the availability of the select.
func (s *MQTTSelect) PublishAvailability(ctx context.Context, msg MQTTMessage) error {
	if len(s.Configuration.Availability) == 0 && s.Configuration.AvailabilityTopic == "" {
		return ErrNoAvailabilityTopic
	}

	topic := ""

	switch {
	case len(s.Configuration.Availability) > 0:
		topic = s.Configuration.Availability[0].Topic
	default:
		topic = s.Configuration.AvailabilityTopic
	}

	s.logger.Debug(
		"publishing availability",
		zap.String("topic", topic),
		zap.String("entity", s.Configuration.Name),
	)

	msgjson, err := msg.Marshal()
	if err != nil {
		return err
	}

	return s.ec.Publish(ctx, topic, msgjson)
}

// PublishState publishes the state of the select.
func (s *MQTTSelect) PublishState(ctx context.Context, msg MQTTMessage) error {
	if s.Configuration.StateTopic == "" {
		return ErrNoStateTopic
	}

	s.logger.Debug(
		"publishing state",
		zap.String("topic", s.Configuration.StateTopic),
		zap.String("entity", s.Configuration.Name),
	)

	msgjson, err := msg.Marshal()
	if err != nil {
		return err
	}

	return s.ec.Publish(ctx, s.Configuration.StateTopic, msgjson)
}

// PublishConfig publishes the configuration of the select.
func (s *MQTTSelect) PublishConfig(ctx context.Context) error {
	s.logger.Debug(
		"publishing config",
		zap.String("topic", s.ConfigTopic),
		zap.String("entity", s.Configuration.Name),
	)

	msgjson, err := json.Marshal(s.Configuration)
	if err != nil {
		return err
	}

	return s.ec.Publish(ctx, s.ConfigTopic, msgjson)
}

// SubscribeCommand subscribes to the command topic.
func (s *MQTTSelect) SubscribeCommand(ctx context.Context) (chan []byte, error) {
	if s.Configuration.CommandTopic == "" {
		return nil, ErrNoCommandTopic
	}

	return s.ec.Subscribe(ctx, s.Configuration.CommandTopic)
}

// UnsubscribeCommand unsubscribes from the command topic.
func (s *MQTTSelect) UnsubscribeCommand(ctx context.Context) error {
	if s.Configuration.CommandTopic == "" {
		return ErrNoCommandTopic
	}

	return s.ec.Unsubscribe(ctx, s.Configuration.CommandTopic)
}
