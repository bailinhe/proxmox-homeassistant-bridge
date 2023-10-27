package vmcontroller

import (
	"context"
	"fmt"

	"gitlab.com/bighotel/proxmox-ha-bridge/pkg/homeassistant"
)

// Select is a select interface
type Select interface {
	// Run starts the select entity
	Run(context.Context)
	// Stop stops the select entity
	Stop()
}

func mkSelectConfig(selectName, selectSlug string, device *homeassistant.MQTTDevice) (
	cfgtopic string, cfg *homeassistant.MQTTSelectConfiguration,
) {
	prefix := fmt.Sprintf("homeassistant/select/%s/%s", device.Name, selectSlug)
	cfgtopic = fmt.Sprintf("%s/config", prefix)
	statetopic := fmt.Sprintf("%s/state", prefix)
	availabilitytopic := fmt.Sprintf("%s/availability", prefix)
	commandTopic := fmt.Sprintf("proxmox-ha-bridge/%s/%s/command", device.Name, selectSlug)

	cfg = &homeassistant.MQTTSelectConfiguration{
		AvailabilityMode: homeassistant.MQTTAvailabilityModeLatest,
		Availability: []homeassistant.MQTTAvailability{
			{
				PayloadAvailable:    string(VMStateAvailabilityOnline),
				PayLoadNotAvailable: string(VMStateAvailabilityOffline),
				Topic:               availabilitytopic,
				ValueTemplate:       "{{ value_json.status }}",
			},
		},
		QoS:             1,
		Device:          device,
		Name:            selectName,
		UniqueID:        fmt.Sprintf("%s/proxmox-ha/%s/%s", device.Manufacturer, device.Name, selectSlug),
		CommandTopic:    commandTopic,
		CommandTemplate: `{ "command": "{{ value }}" }`,
		Options: []string{
			string(VMCommandOptionStart),
			string(VMCommandOptionStop),
			string(VMCommandOptionShutdown),
			string(VMCommandOptionReboot),
			string(VMCommandOptionReset),
		},
		StateTopic:             statetopic,
		JSONAttributesTopic:    statetopic,
		JSONAttributesTemplate: `{{ value_json | tojson }}`,
		ValueTemplate:          "{{ value_json.state }}",
	}

	return cfgtopic, cfg
}
