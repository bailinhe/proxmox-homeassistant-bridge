package proxmox

import "errors"

// ErrListingVM is an error that occurs when listing VMs
var ErrListingVM = errors.New("error listing vms")

// ErrTaskFailed is an error that occurs when a task fails
var ErrTaskFailed = errors.New("task failed")
