package homeassistant

// MQTTEventConfiguration is the MQTT Event Config type.
// The mqtt event platform allows you to process event info from an MQTT
// message. Events are signals that are emitted when something happens, for
// example, when a user presses a physical button like a doorbell or when a
// button on a remote control is pressed. With the event some event attributes
// can be sent te become available as an attribute on the entity. MQTT events
// are stateless. For example, a doorbell does not have a state like being “on”
// or “off” but instead is momentarily pressed.
type MQTTEventConfiguration struct {
	// availability list (optional)
	// A list of MQTT topics subscribed to receive availability (online/offline)
	// updates. Must not be used together with availability_topic.
	Avalability []MQTTAvailability `json:"availability,omitempty"`

	// availability_mode string (optional, default: latest)
	// When availability is configured, this controls the conditions needed to set
	// the entity to available. Valid entries are all, any, and latest. If set to
	// all, payload_available must be received on all configured availability
	// topics before the entity is marked as online. If set to any, payload_available
	// must be received on at least one configured availability topic before the
	// entity is marked as online. If set to latest, the last payload_available or
	// payload_not_available received on any configured availability topic controls
	// the availability.
	AvalabilityMode MQTTAvailabilityMode `json:"availability_mode,omitempty"`

	// availability_template template (optional)
	// Defines a template to extract device’s availability from the availability_topic.
	// To determine the devices’s availability result of this template will be compared
	// to payload_available and payload_not_available.
	AvalabilityTemplate string `json:"availability_template,omitempty"`

	// availability_topic string (optional)
	// The MQTT topic subscribed to receive availability (online/offline) updates.
	// Must not be used together with availability.
	AvalabilityTopic string `json:"availability_topic,omitempty"`

	// device map (optional)
	// Information about the device this event is a part of to tie it into the device
	// registry. Only works when unique_id is set. At least one of identifiers or
	// connections must be present to identify the device.
	Device *MQTTDevice `json:"device,omitempty"`

	// device_class device_class (optional, default: None)
	// The type/class of the event to set the icon in the frontend. The device_class
	// can be null.
	DeviceClass string `json:"device_class,omitempty"`

	// enabled_by_default boolean (optional, default: true)
	// Flag which defines if the entity should be enabled when first added.
	EneabledByDefault bool `json:"enabled_by_default,omitempty"`

	// encoding string (optional, default: utf-8)
	// The encoding of the published messages.
	Eneabled string `json:"encoding,omitempty"`

	// entity_category string (optional, default: None)
	// The category of the entity.
	EntityCategory string `json:"entity_category,omitempty"`

	// event_types list REQUIRED
	// A list of valid event_type strings.
	EventTypes []string `json:"event_types"`

	// icon icon (optional)
	// Icon for the entity.
	Icon string `json:"icon,omitempty"`

	// json_attributes_template template (optional)
	// Defines a template to extract the JSON dictionary from messages received on
	// the json_attributes_topic. Usage example can be found in MQTT sensor documentation.
	JSONAttributesTemplate string `json:"json_attributes_template,omitempty"`

	// json_attributes_topic string (optional)
	// The MQTT topic subscribed to receive a JSON dictionary payload and then set
	// as sensor attributes. Usage example can be found in MQTT sensor documentation.
	JSONAttributesTopic string `json:"json_attributes_topic,omitempty"`

	// name string (optional, default: MQTT Event)
	// The name to use when displaying this event.
	Name string `json:"name,omitempty"`

	// object_id string (optional)
	// Used instead of name for automatic generation of entity_id
	ObjectID string `json:"object_id,omitempty"`

	// payload_available string (optional, default: online)
	// The payload that represents the available state.
	PayloadAvailable string `json:"payload_available,omitempty"`

	// payload_not_available string (optional, default: offline)
	// The payload that represents the unavailable state.
	PayloadNotAvailable string `json:"payload_not_available,omitempty"`

	// qos integer (optional, default: 0)
	// The maximum QoS level to be used when receiving and publishing messages.
	QoS int `json:"qos,omitempty"`

	// state_topic string REQUIRED, default: None
	// The MQTT topic subscribed to receive JSON event payloads. The JSON payload
	// should contain the event_type element. The event type should be one of the
	// configured event_types.
	StateTopic string `json:"state_topic"`

	// unique_id string (optional)
	// An ID that uniquely identifies this event entity. If two events have the same
	// unique ID, Home Assistant will raise an exception.
	UniqueID string `json:"unique_id,omitempty"`

	// value_template template (optional)
	// Defines a template to extract the value and render it to a valid JSON event
	// payload. If the template throws an error, the current state will be used instead.
	ValueTemplate string `json:"value_template,omitempty"`
}
