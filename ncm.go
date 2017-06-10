package gosolar

import "fmt"

// RemoveNCMNodes is awesome.
func (c *Client) RemoveNCMNodes(guids []string) error {
	req := [][]string{guids}

	_, err := c.post("Invoke/Cirrus.Nodes/RemoveNodes", req)
	if err != nil {
		return fmt.Errorf("failed to remove the NCM nodes %v", err)
	}

	return nil
}
