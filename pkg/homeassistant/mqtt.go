package homeassistant

// MQTTAvailability is the configuration for the MQTT availability
type MQTTAvailability struct {
	PayloadAvailable    string `json:"payload_available"`
	PayLoadNotAvailable string `json:"payload_not_available"`
	Topic               string `json:"topic"`
	ValueTemplate       string `json:"value_template,omitempty"`
}

// MQTTAvailabilityMode is an enum for the MQTT availability mode
// When availability is configured, this controls the conditions needed to set
// the entity to available. Valid entries are all, any, and latest. If set to
// all, payload_available must be received on all configured availability
// topics before the entity is marked as online. If set to any, payload_available
// must be received on at least one configured availability topic before the
// entity is marked as online. If set to latest, the last payload_available or
// payload_not_available received on any configured availability topic controls
// the availability.
type MQTTAvailabilityMode string

const (
	// MQTTAvailabilityModeAny is the any availability mode
	MQTTAvailabilityModeAny MQTTAvailabilityMode = "any"
	// MQTTAvailabilityModeAll is the all availability mode
	MQTTAvailabilityModeAll MQTTAvailabilityMode = "all"
	// MQTTAvailabilityModeLatest is the latest availability mode
	MQTTAvailabilityModeLatest MQTTAvailabilityMode = "latest"
)

// MQTTDevice is the configuration for the MQTT device
type MQTTDevice struct {
	ConfigurationURL string     `json:"configuration_url,omitempty"`
	Connection       [][]string `json:"connections,omitempty"`
	HardwareVersion  string     `json:"hw_version,omitempty"`
	Identifiers      []string   `json:"identifiers,omitempty"`
	Manufacturer     string     `json:"manufacturer,omitempty"`
	Model            string     `json:"model,omitempty"`
	Name             string     `json:"name,omitempty"`
	SuggestedArea    string     `json:"suggested_area,omitempty"`
	SoftwareVersion  string     `json:"sw_version,omitempty"`
	ViaDevice        string     `json:"via_device,omitempty"`
}
