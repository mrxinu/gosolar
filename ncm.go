package gosolar

import "fmt"

// RemoveNCMNodes is now even more awesome.
func (c *Client) RemoveNCMNodes(guids []string) error {
	req, endpoint := getRemoveNCMNodesRequest(guids)
	_, err := c.post(endpoint, req)

	if err != nil {
		return fmt.Errorf("failed to remove the NCM nodes %v", err)
	}

	return nil
}

// RemoveNodeEndpoint is the endpoint to send the post request to remove NCM Nodes
const RemoveNodeEndpoint = "Invoke/Cirrus.Nodes/RemoveNodes"

// getRemoveNCMNodesRequest is a function that will convert a slice of guid strings into
// an endpoint and a request that the API expects.
func getRemoveNCMNodesRequest(guids []string) ([][]string, string) {
	req := [][]string{guids}
	return req, RemoveNodeEndpoint
}
