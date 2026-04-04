package gosolar

import "context"

// RemoveNCMNodes deletes nodes from NCM handling in SolarWinds.
func (c *Client) RemoveNCMNodes(ctx context.Context, guids []string) error {
	_, err := c.Invoke(ctx, "Cirrus.Nodes", "RemoveNodes", [][]string{guids})
	if err != nil {
		return WrapError(err, ErrorTypeInternal, "remove_ncm_nodes", "failed to remove NCM nodes")
	}
	return nil
}
