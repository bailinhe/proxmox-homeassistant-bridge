package proxmox

import "github.com/luthermonson/go-proxmox"

// Nodes returns the nodes from the Proxmox API
func (c *Client) Nodes() (proxmox.NodeStatuses, error) {
	if c.nodes != nil {
		return c.nodes, nil
	}

	nodes, err := c.client.Nodes()
	if err != nil {
		return nil, err
	}

	c.nodes = nodes

	return c.nodes, nil
}
