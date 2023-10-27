package proxmox

import (
	"context"
	"fmt"
	"time"

	goproxmox "github.com/luthermonson/go-proxmox"
	"gitlab.com/bighotel/proxmox-ha-bridge/pkg/vmcontroller"
	"go.uber.org/zap"
)

// VM is a wrapper for the proxmox VM
type VM struct {
	vm     *goproxmox.VirtualMachine
	client *goproxmox.Client
	node   *goproxmox.Node
	logger *zap.Logger
}

// NewVM creates a new proxmox VM controller
func NewVM(vm *goproxmox.VirtualMachine, node *goproxmox.Node, client *goproxmox.Client, logger *zap.Logger) *VM {
	vmctl := &VM{
		vm:     vm,
		client: client,
		logger: logger,
		node:   node,
	}

	vmctl.logger = vmctl.logger.With(zap.String("component", "proxmox-vm-controller"))

	return vmctl
}

// VM implements VMController interface
var _ vmcontroller.VMController = (*VM)(nil)

func (p *VM) waitTask(ctx context.Context, task *goproxmox.Task) error {
	// wait for the task to complete
	maxWaitMinutes := 2

	p.logger.Debug(
		"waiting for task to complete",
		zap.Int64("vmid", int64(p.vm.VMID)),
		zap.String("task", string(task.UPID)),
		zap.Int("max-wait-minutes", maxWaitMinutes),
	)

	errChan := make(chan error)

	go func() {
		if err := task.Wait(time.Second, time.Duration(maxWaitMinutes)*time.Minute); err != nil {
			errChan <- err
		}

		if !task.IsSuccessful {
			errChan <- fmt.Errorf("%w: %s", ErrTaskFailed, task.ExitStatus)
		}

		errChan <- nil
	}()

	select {
	case <-ctx.Done():
		p.logger.Debug("task wait cancelled", zap.Int64("vmid", int64(p.vm.VMID)), zap.String("task", string(task.UPID)))
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}

// Start starts the vm
func (p *VM) Start(ctx context.Context) error {
	p.logger.Info("starting vm", zap.Int64("vmid", int64(p.vm.VMID)))

	vmstatus, err := p.Status()
	if err != nil {
		p.logger.Error("error getting vm status", zap.Error(err))
		return err
	}

	if vmstatus == "running" {
		p.logger.Warn("vm is already running", zap.Int64("vmid", int64(p.vm.VMID)))
		return nil
	}

	startvm := func() error {
		task, err := p.vm.Start()
		if err != nil {
			return err
		}

		return p.waitTask(ctx, task)
	}

	// five times and then give up
	maxTries := 5

	for i := 0; i < maxTries; i++ {
		if err = startvm(); err != nil {
			p.logger.Error("error starting vm", zap.Error(err))
			p.logger.Sugar().Infof("retrying in %d minutes", i)

			time.Sleep(time.Duration(i) * time.Second)

			continue
		}

		break
	}

	if err == nil {
		p.logger.Info("vm started", zap.Int64("vmid", int64(p.vm.VMID)))
	}

	return err
}

// Stop stops the vm
func (p *VM) Stop(ctx context.Context) error {
	p.logger.Info("stopping vm", zap.Int64("vmid", int64(p.vm.VMID)))

	vmstatus, err := p.Status()
	if err != nil {
		p.logger.Error("error getting vm status", zap.Error(err))
		return err
	}

	if vmstatus == "stopped" {
		p.logger.Warn("vm is already stopped", zap.Int64("vmid", int64(p.vm.VMID)))
		return nil
	}

	stopvm := func() error {
		task, err := p.vm.Stop()
		if err != nil {
			return err
		}

		return p.waitTask(ctx, task)
	}

	return stopvm()
}

// Restart restarts the vm
func (p *VM) Restart(ctx context.Context) error {
	p.logger.Info("restarting vm", zap.Int64("vmid", int64(p.vm.VMID)))

	vmstatus, err := p.Status()
	if err != nil {
		p.logger.Error("error getting vm status", zap.Error(err))
		return err
	}

	if vmstatus == "stopped" {
		p.logger.Warn("vm is not running", zap.Int64("vmid", int64(p.vm.VMID)))
		return nil
	}

	restartvm := func() error {
		task, err := p.vm.Reboot()
		if err != nil {
			return err
		}

		return p.waitTask(ctx, task)
	}

	return restartvm()
}

// Reset hard resets the vm
func (p *VM) Reset(ctx context.Context) error {
	p.logger.Info("resetting vm", zap.Int64("vmid", int64(p.vm.VMID)))

	vmstatus, err := p.Status()
	if err != nil {
		p.logger.Error("error getting vm status", zap.Error(err))
		return err
	}

	if vmstatus == "stopped" {
		p.logger.Warn("vm is not running", zap.Int64("vmid", int64(p.vm.VMID)))
		return nil
	}

	resetvm := func() error {
		task, err := p.vm.Reset()
		if err != nil {
			return err
		}

		return p.waitTask(ctx, task)
	}

	return resetvm()
}

// Shutdown gracefully shuts down the vm
func (p *VM) Shutdown(ctx context.Context) error {
	p.logger.Info("shutting down vm", zap.Int64("vmid", int64(p.vm.VMID)))

	vmstatus, err := p.Status()
	if err != nil {
		p.logger.Error("error getting vm status", zap.Error(err))
		return err
	}

	if vmstatus == "stopped" {
		p.logger.Warn("vm is not running", zap.Int64("vmid", int64(p.vm.VMID)))
		return nil
	}

	shutdownvm := func() error {
		task, err := p.vm.Shutdown()
		if err != nil {
			return err
		}

		return p.waitTask(ctx, task)
	}

	return shutdownvm()
}

// Status returns the status of the vm
func (p *VM) Status() (string, error) {
	p.logger.Debug("getting status of vm", zap.Int64("vmid", int64(p.vm.VMID)))

	vm, err := p.node.VirtualMachine(int(p.vm.VMID))
	if err != nil {
		return "", err
	}

	p.vm = vm

	return p.vm.Status, nil
}

// ID returns the id of the vm
func (p *VM) ID() string {
	return fmt.Sprintf("%d", p.vm.VMID)
}

// Name returns the name of the vm
func (p *VM) Name() string {
	return p.vm.Name
}
