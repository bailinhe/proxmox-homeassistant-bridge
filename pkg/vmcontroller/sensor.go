package vmcontroller

import (
	"context"
	"fmt"

	"gitlab.com/bighotel/proxmox-ha-bridge/pkg/homeassistant"
)

// Sensor is a sensor interface
type Sensor interface {
	// Run starts the sensor
	Run(context.Context)
	// Stop stops the sensor
	Stop()
}

func mkSensorConfig(sensorName, sensorSlug string, device *homeassistant.MQTTDevice) (
	cfgtopic string, cfg *homeassistant.MQTTSensorConfiguration,
) {
	prefix := fmt.Sprintf("homeassistant/sensor/%s/%s", device.Name, sensorSlug)
	cfgtopic = fmt.Sprintf("%s/config", prefix)
	statetopic := fmt.Sprintf("%s/state", prefix)
	availabilitytopic := fmt.Sprintf("%s/availability", prefix)

	cfg = &homeassistant.MQTTSensorConfiguration{
		AvailabilityMode: homeassistant.MQTTAvailabilityModeLatest,
		Availability: []homeassistant.MQTTAvailability{
			{
				PayloadAvailable:    string(VMStateAvailabilityOnline),
				PayLoadNotAvailable: string(VMStateAvailabilityOffline),
				Topic:               availabilitytopic,
				ValueTemplate:       "{{ value_json.status }}",
			},
		},
		Device:                 device,
		Name:                   sensorName,
		UniqueID:               fmt.Sprintf("%s/proxmox-ha/%s/%s", device.Manufacturer, device.Name, sensorSlug),
		JSONAttributesTemplate: `{{ value_json | tojson }}`,
		JSONAttributesTopic:    statetopic,
		StateTopic:             statetopic,
		ValueTemplate:          "{{ value_json.state }}",
	}

	return
}
