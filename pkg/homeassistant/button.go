package homeassistant

// MQTTButtonConfiguration is the configuration for a MQTT button
type MQTTButtonConfiguration struct {
	Availability           []MQTTAvailability   `json:"availability,omitempty"`
	PayloadAvailable       string               `json:"payload_available,omitempty" default:"online"`
	PayloadNotAvailable    string               `json:"payload_not_available,omitempty" default:"offline"`
	Topic                  string               `json:"topic"`
	ValueTemplate          string               `json:"value_template,omitempty"`
	AvailabilityMode       MQTTAvailabilityMode `json:"availability_mode,omitempty" default:"latest"`
	AvailabilityTemplate   string               `json:"availability_template,omitempty"`
	AvailabilityTopic      string               `json:"availability_topic,omitempty"`
	CommandTemplate        string               `json:"command_template,omitempty"`
	CommandTopic           string               `json:"command_topic,omitempty"`
	Device                 *MQTTDevice          `json:"device,omitempty"`
	DeviceClass            string               `json:"device_class,omitempty" default:"None"`
	EnabledByDefault       bool                 `json:"enabled_by_default,omitempty" default:"true"`
	Encoding               string               `json:"encoding,omitempty" default:"utf-8"`
	EntityCategory         string               `json:"entity_category,omitempty" default:"None"`
	Icon                   string               `json:"icon,omitempty"`
	JSONAttributesTemplate string               `json:"json_attributes_template,omitempty"`
	JSONAttributesTopic    string               `json:"json_attributes_topic,omitempty"`
	ObjectID               string               `json:"object_id,omitempty"`
	PayloadPress           string               `json:"payload_press,omitempty" default:"PRESS"`
	QoS                    int                  `json:"qos,omitempty" default:"0"`
	Retain                 bool                 `json:"retain,omitempty" default:"false"`
	UniqueID               string               `json:"unique_id,omitempty"`
}
