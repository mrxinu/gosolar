package gosolar

import (
	"context"
	"encoding/json"
)

// Assignment holds a Universal Device Poller (UnDP) assignment from SolarWinds.
type Assignment struct {
	ID             string `json:"CustomPollerAssignmentID"`
	PollerID       string `json:"PollerID"`
	NodeID         int    `json:"NodeID"`
	InterfaceID    int    `json:"InterfaceID"`
	CustomPollerID string `json:"CustomPollerID"`
	InstanceType   string `json:"InstanceType"`
}

// GetAssignments returns all current custom poller assignments.
func (c *Client) GetAssignments(ctx context.Context) ([]Assignment, error) {
	query := `
		SELECT
			CustomPollerAssignmentID,
			CustomPollerID,
			NodeID,
			InterfaceID,
			CustomPollerID,
			InstanceType
		FROM Orion.NPM.CustomPollerAssignment
	`

	res, err := c.Query(ctx, query, nil)
	if err != nil {
		return nil, WrapError(err, ErrorTypeInternal, "get_assignments", "failed to query for assignments")
	}

	var assignments []Assignment
	if err := json.Unmarshal(res, &assignments); err != nil {
		return nil, WrapError(err, ErrorTypeInternal, "get_assignments", "failed to unmarshal assignments")
	}

	return assignments, nil
}

// AddNodePoller adds a Universal Device Poller (UnDP) to a node.
func (c *Client) AddNodePoller(ctx context.Context, customPollerID string, nodeID int) error {
	request := struct {
		NodeID         int    `json:"NodeID"`
		CustomPollerID string `json:"CustomPollerID"`
	}{
		NodeID:         nodeID,
		CustomPollerID: customPollerID,
	}

	_, err := c.Create(ctx, "Orion.NPM.CustomPollerAssignmentOnNode", request)
	if err != nil {
		return WrapError(err, ErrorTypeInternal, "add_node_poller", "failed to add poller")
	}

	return nil
}

// AddInterfacePoller adds a Universal Device Poller (UnDP) to an interface.
func (c *Client) AddInterfacePoller(ctx context.Context, customPollerID string, interfaceID int) error {
	request := struct {
		InterfaceID    int    `json:"InterfaceID"`
		CustomPollerID string `json:"CustomPollerID"`
	}{
		InterfaceID:    interfaceID,
		CustomPollerID: customPollerID,
	}

	_, err := c.Create(ctx, "Orion.NPM.CustomPollerAssignmentOnInterface", request)
	if err != nil {
		return WrapError(err, ErrorTypeInternal, "add_interface_poller", "failed to add poller")
	}

	return nil
}
