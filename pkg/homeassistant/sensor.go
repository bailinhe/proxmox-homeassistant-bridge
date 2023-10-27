package homeassistant

import (
	"context"
	"encoding/json"

	"gitlab.com/bighotel/proxmox-ha-bridge/pkg/events"
	"go.uber.org/zap"
)

// SensorDeviceClass is an enum type for the device class of a sensor
type SensorDeviceClass string

const (
	// SensorDeviceClassNone is the none device class
	SensorDeviceClassNone SensorDeviceClass = "None"
	// SensorDeviceClassApparentPower is the apparent power device class
	SensorDeviceClassApparentPower SensorDeviceClass = "apparent_power"
	// SensorDeviceClassAqi is the air quality index device class
	SensorDeviceClassAqi SensorDeviceClass = "aqi"
	// SensorDeviceClassBattery is the battery device class
	SensorDeviceClassBattery SensorDeviceClass = "battery"
	// SensorDeviceClassCarbonDioxide is the carbon dioxide device class
	SensorDeviceClassCarbonDioxide SensorDeviceClass = "carbon_dioxide"
	// SensorDeviceClassCarbonMonoxide is the carbon monoxide device class
	SensorDeviceClassCarbonMonoxide SensorDeviceClass = "carbon_monoxide"
	// SensorDeviceClassCurrent is the current device class
	SensorDeviceClassCurrent SensorDeviceClass = "current"
	// SensorDeviceClassDate is the date device class
	SensorDeviceClassDate SensorDeviceClass = "date"
	// SensorDeviceClassDuration is the duration device class
	SensorDeviceClassDuration SensorDeviceClass = "duration"
	// SensorDeviceClassEnergy is the energy device class
	SensorDeviceClassEnergy SensorDeviceClass = "energy"
	// SensorDeviceClassFrequency is the frequency device class
	SensorDeviceClassFrequency SensorDeviceClass = "frequency"
	// SensorDeviceClassGas is the gas device class
	SensorDeviceClassGas SensorDeviceClass = "gas"
	// SensorDeviceClassHumidity is the humidity device class
	SensorDeviceClassHumidity SensorDeviceClass = "humidity"
	// SensorDeviceClassIlluminance is the illuminance device class
	SensorDeviceClassIlluminance SensorDeviceClass = "illuminance"
	// SensorDeviceClassMonetary is the monetary device class
	SensorDeviceClassMonetary SensorDeviceClass = "monetary"
	// SensorDeviceClassNitrogenDioxide is the nitrogen dioxide device class
	SensorDeviceClassNitrogenDioxide SensorDeviceClass = "nitrogen_dioxide"
	// SensorDeviceClassNitrogenMonoxide is the nitrogen monoxide device class
	SensorDeviceClassNitrogenMonoxide SensorDeviceClass = "nitrogen_monoxide"
	// SensorDeviceClassNitrousOxide is the nitrous oxide device class
	SensorDeviceClassNitrousOxide SensorDeviceClass = "nitrous_oxide"
	// SensorDeviceClassOzone is the ozone device class
	SensorDeviceClassOzone SensorDeviceClass = "ozone"
	// SensorDeviceClassPm1 is the PM1 device class
	SensorDeviceClassPm1 SensorDeviceClass = "pm1"
	// SensorDeviceClassPm10 is the PM10 device class
	SensorDeviceClassPm10 SensorDeviceClass = "pm10"
	// SensorDeviceClassPm25 is the PM2.5 device class
	SensorDeviceClassPm25 SensorDeviceClass = "pm25"
	// SensorDeviceClassPowerFactor is the power factor device class
	SensorDeviceClassPowerFactor SensorDeviceClass = "power_factor"
	// SensorDeviceClassPower is the power device class
	SensorDeviceClassPower SensorDeviceClass = "power"
	// SensorDeviceClassPressure is the pressure device class
	SensorDeviceClassPressure SensorDeviceClass = "pressure"
	// SensorDeviceClassReactivePower is the reactive power device class
	SensorDeviceClassReactivePower SensorDeviceClass = "reactive_power"
	// SensorDeviceClassSignalStrength is the signal strength device class
	SensorDeviceClassSignalStrength SensorDeviceClass = "signal_strength"
	// SensorDeviceClassSulphurDioxide is the sulphur dioxide device class
	SensorDeviceClassSulphurDioxide SensorDeviceClass = "sulphur_dioxide"
	// SensorDeviceClassTemperature is the temperature device class
	SensorDeviceClassTemperature SensorDeviceClass = "temperature"
	// SensorDeviceClassTimestamp is the timestamp device class
	SensorDeviceClassTimestamp SensorDeviceClass = "timestamp"
	// SensorDeviceClassVolatileOrganicCompounds is the volatile organic compounds device class
	SensorDeviceClassVolatileOrganicCompounds SensorDeviceClass = "volatile_organic_compounds"
	// SensorDeviceClassVoltage is the voltage device class
	SensorDeviceClassVoltage SensorDeviceClass = "voltage"
)

// SensorStateClass is an enum type for the state class of a sensor
type SensorStateClass string

const (
	// SensorStateClassNone is the none state class
	SensorStateClassNone SensorStateClass = "None"
	// SensorStateClassMeasurement is the measurement state class
	SensorStateClassMeasurement SensorStateClass = "measurement"
	// SensorStateClassTotal is the total state class
	SensorStateClassTotal SensorStateClass = "total"
	// SensorStateClassTotalIncreasing is the total increasing state class
	SensorStateClassTotalIncreasing SensorStateClass = "total_increasing"
)

// MQTTSensorConfiguration is the configuration for a MQTT sensor
type MQTTSensorConfiguration struct {
	Availability           []MQTTAvailability   `json:"availability,omitempty"`
	AvailabilityMode       MQTTAvailabilityMode `json:"availability_mode,omitempty" default:"latest"`
	AvailabilityTemplate   string               `json:"availability_template,omitempty"`
	AvailabilityTopic      string               `json:"availability_topic,omitempty"`
	Device                 *MQTTDevice          `json:"device,omitempty"`
	DeviceClass            SensorDeviceClass    `json:"device_class,omitempty" default:"None"`
	EnabledByDefault       bool                 `json:"enabled_by_default,omitempty" default:"true"`
	Encoding               string               `json:"encoding,omitempty" default:"utf-8"`
	EntityCategory         string               `json:"entity_category,omitempty" default:"None"`
	ExpireAfter            int                  `json:"expire_after,omitempty" default:"0"`
	ForceUpdate            bool                 `json:"force_update,omitempty" default:"false"`
	Icon                   string               `json:"icon,omitempty"`
	JSONAttributesTemplate string               `json:"json_attributes_template,omitempty"`
	JSONAttributesTopic    string               `json:"json_attributes_topic,omitempty"`
	LastResetValueTemplate string               `json:"last_reset_value_template,omitempty"`
	Name                   string               `json:"name,omitempty" default:"MQTT Sensor"`
	ObjectID               string               `json:"object_id,omitempty"`
	PayloadAvailable       string               `json:"payload_available"`
	PayLoadNotAvailable    string               `json:"payload_not_available"`
	QoS                    int                  `json:"qos:omitempty" default:"0"`
	StateClass             SensorStateClass     `json:"state_class,omitempty" default:"None"`
	StateTopic             string               `json:"state_topic"`
	UniqueID               string               `json:"unique_id,omitempty"`
	UnitOfMeasurement      string               `json:"unit_of_measurement,omitempty"`
	ValueTemplate          string               `json:"value_template,omitempty"`
}

// MQTTMessage is a MQTT message
type MQTTMessage interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

// MQTTSensor is a MQTT sensor
type MQTTSensor struct {
	Configuration *MQTTSensorConfiguration
	ConfigTopic   string
	ec            events.Client
	logger        *zap.Logger
}

// MQTTSensorOpt is a functional option for the MQTTSensor
type MQTTSensorOpt func(*MQTTSensor)

// MQTTSensorWithLogger sets the logger for the MQTTSensor
func MQTTSensorWithLogger(logger *zap.Logger) MQTTSensorOpt {
	return func(s *MQTTSensor) {
		s.logger = logger
	}
}

// NewMQTTSensor creates a new MQTTSensor
func NewMQTTSensor(cfg *MQTTSensorConfiguration, cfgtopic string, ec events.Client, opts ...MQTTSensorOpt) *MQTTSensor {
	s := &MQTTSensor{
		Configuration: cfg,
		ConfigTopic:   cfgtopic,
		ec:            ec,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.logger = s.logger.With(zap.String("component", "mqtt-sensor"))

	return s
}

// PublishAvailability publishes the availability of the sensor
func (s *MQTTSensor) PublishAvailability(ctx context.Context, msg MQTTMessage) error {
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

// PublishState publishes the state of the sensor
func (s *MQTTSensor) PublishState(ctx context.Context, msg MQTTMessage) error {
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

// PublishConfig publishes the configuration of the sensor
func (s *MQTTSensor) PublishConfig(ctx context.Context) error {
	cfgjson, err := json.Marshal(s.Configuration)
	if err != nil {
		return err
	}

	s.logger.Debug(
		"publishing config",
		zap.String("topic", s.ConfigTopic),
		zap.String("entity", s.Configuration.Name),
	)

	return s.ec.Publish(ctx, s.ConfigTopic, cfgjson)
}
