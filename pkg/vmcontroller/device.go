package vmcontroller

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/bighotel/proxmox-ha-bridge/pkg/events"
	"gitlab.com/bighotel/proxmox-ha-bridge/pkg/homeassistant"
	"go.uber.org/zap"
)

const (
	// DefaultProbeInterval is the default probe interval for the VMStateSensor
	DefaultProbeInterval = 1 * time.Minute
	// DefaultAvailabilityPublishInterval is the default availability publish
	// interval for the VMStateSensor
	DefaultAvailabilityPublishInterval = 1 * time.Minute
)

// DeviceTemplate is the MQTT device template for this app
var DeviceTemplate = homeassistant.MQTTDevice{
	Manufacturer:    "bailinhe.com",
	Model:           "Proxmox Home Assistant Bridge",
	Name:            "proxmox-ha-bridge",
	HardwareVersion: "1.0.0",
	SoftwareVersion: "0.0.1",
}

// HomeAssistantEntity is an interface for a Home Assistant entity
type HomeAssistantEntity interface {
	// Run starts the entity
	Run(context.Context)
	// Stop stops the entity
	Stop()
}

// VirtualMachineMQTTDevice is a virtual machine MQTT device
type VirtualMachineMQTTDevice struct {
	vm VMController
	ec events.Client

	entities                    []HomeAssistantEntity
	logger                      *zap.Logger
	probeInterval               time.Duration
	availabilityPublishInterval time.Duration

	stopchan chan struct{}
}

// MQTTDeviceOpt is a virtual machine MQTT device option
type MQTTDeviceOpt func(*VirtualMachineMQTTDevice)

// WithLogger sets the logger
func WithLogger(logger *zap.Logger) MQTTDeviceOpt {
	return func(d *VirtualMachineMQTTDevice) {
		d.logger = logger
	}
}

// WithProbeInterval sets the probe interval
func WithProbeInterval(t time.Duration) MQTTDeviceOpt {
	return func(d *VirtualMachineMQTTDevice) {
		d.probeInterval = t
	}
}

// WithAvailabilityPublishInterval sets the availability publish interval
func WithAvailabilityPublishInterval(t time.Duration) MQTTDeviceOpt {
	return func(d *VirtualMachineMQTTDevice) {
		d.availabilityPublishInterval = t
	}
}

// NewVirtualMachineMQTTDevice creates a new virtual machine MQTT device
func NewVirtualMachineMQTTDevice(vm VMController, ec events.Client, opts ...MQTTDeviceOpt) *VirtualMachineMQTTDevice {
	vmd := &VirtualMachineMQTTDevice{
		vm:                          vm,
		ec:                          ec,
		probeInterval:               DefaultProbeInterval,
		availabilityPublishInterval: DefaultAvailabilityPublishInterval,
		stopchan:                    make(chan struct{}),
	}

	for _, opt := range opts {
		opt(vmd)
	}

	vmd.logger = vmd.logger.With(zap.String("component", "vm-mqtt-device"))

	return vmd
}

// Run runs the virtual machine MQTT device
func (d *VirtualMachineMQTTDevice) Run(ctx context.Context) {
	d.logger.Info("starting vm mqtt device", zap.String("vm", d.vm.Name()))

	device := &homeassistant.MQTTDevice{
		Identifiers: []string{
			fmt.Sprintf(
				"proxmox-ha-bridge.int.bailinhe.com/vms/%s-%s",
				d.vm.ID(), d.vm.Name(),
			),
		},
		Manufacturer: DeviceTemplate.Manufacturer,
		Model:        DeviceTemplate.Model,
		Name:         fmt.Sprintf("vm-%s-%s", d.vm.ID(), d.vm.Name()),
	}

	d.entities = []HomeAssistantEntity{
		NewVMStateSensor(
			d.vm, d.ec, device,
			VMStateSensorWithLogger(d.logger),
			VMStateSensorWithAvailabilityPublishInterval(d.availabilityPublishInterval),
			VMStateSensorWithProbeInterval(d.probeInterval),
		),

		NewVMCommandSelect(
			d.vm, d.ec, device,
			VMCommandSelectWithLogger(d.logger),
			VMCommandSelectWithAvailabilityPublishInterval(d.availabilityPublishInterval),
		),
	}

	for _, s := range d.entities {
		go s.Run(ctx)
	}

	<-d.stopchan

	for _, s := range d.entities {
		s.Stop()
	}
}

// Stop stops the virtual machine MQTT device
func (d *VirtualMachineMQTTDevice) Stop() {
	d.logger.Info("stopping vm mqtt device", zap.String("vm", d.vm.Name()))
	d.stopchan <- struct{}{}
}
