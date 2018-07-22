package gosolar

import "fmt"

// RemoveNCMNodes deletes nodes from NCM handling in SolarWinds.
func (c *Client) RemoveNCMNodes(guids []string) error {
	endpoint := "Invoke/Cirrus.Nodes/RemoveNodes"
	req := [][]string{guids}

	_, err := c.post(endpoint, req)

	if err != nil {
		return fmt.Errorf("failed to remove the NCM nodes %v", err)
	}

	return nil
}
