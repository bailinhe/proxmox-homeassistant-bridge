package proxmox

import (
	"fmt"
)

// VMs returns the VMs from the Proxmox API
func (c *Client) VMs() (vms []*VM, errListVMs error) {
	nodes, err := c.Nodes()
	if err != nil {
		errListVMs = fmt.Errorf("%w: %s", ErrListingVM, errListVMs.Error())
		return
	}

	for _, st := range nodes {
		node, err := c.client.Node(st.Node)
		if err != nil {
			errListVMs = fmt.Errorf("%w: %s", ErrListingVM, errListVMs.Error())
			return
		}

		apivms, err := node.VirtualMachines()
		if err != nil {
			errListVMs = fmt.Errorf("%w: %s", ErrListingVM, errListVMs.Error())
			return
		}

		for _, vm := range apivms {
			vms = append(vms, NewVM(vm, node, c.client, c.logger))
		}
	}

	return
}
