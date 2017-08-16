package gosolar

import (
	"encoding/json"
	"fmt"
)

// Assignment holds all the current UnDP configuration from SolarWinds.
type Assignment struct {
	ID             string `json:"CustomPollerAssignmentID"`
	PollerID       string `json:"PollerID"`
	NodeID         int    `json:"NodeID"`
	InterfaceID    int    `json:"InterfaceID"`
	CustomPollerID string `json:"CustomPollerID"`
	InstanceType   string `json:"InstanceType"`
}

// GetAssignments function returns all the current custom poller assignments
// in effect at the time.
func (c *Client) GetAssignments() ([]Assignment, error) {
	query := `
		SELECT
			CustomPollerAssignmentID
			,CustomPollerID
			,NodeID
			,InterfaceID
			,CustomPollerID
			,InstanceType
		FROM Orion.NPM.CustomPollerAssignment
	`

	res, err := c.Query(query, nil)
	if err != nil {
		return []Assignment{}, fmt.Errorf("failed to query for assignments: %v", err)
	}

	var assignments []Assignment
	if err := json.Unmarshal(res, &assignments); err != nil {
		return []Assignment{}, fmt.Errorf("failed to unmarshal assignments: %v", err)
	}

	return assignments, nil
}

// AddNodePoller adds a Universal Device Poller (UnDP) to a node.
func (c *Client) AddNodePoller(customPollerID string, nodeID int) error {
	entity := "Orion.NPM.CustomPollerAssignmentOnNode"

	request := struct {
		NodeID         int    `json:"NodeID"`
		CustomPollerID string `json:"CustomPollerID"`
	}{
		NodeID:         nodeID,
		CustomPollerID: customPollerID,
	}

	_, err := c.post("Create/"+entity, request)
	if err != nil {
		return fmt.Errorf("failed to add poller: %v", err)
	}

	return nil
}

// AddInterfacePoller adds a Universal Device Poller (UnDP) to an interface.
func (c *Client) AddInterfacePoller(customPollerID string, interfaceID int) error {
	entity := "Orion.NPM.CustomPollerAssignmentOnInterface"

	request := struct {
		InterfaceID    int    `json:"InterfaceID"`
		CustomPollerID string `json:"CustomPollerID"`
	}{
		InterfaceID:    interfaceID,
		CustomPollerID: customPollerID,
	}

	_, err := c.post("Create/"+entity, request)
	if err != nil {
		return fmt.Errorf("failed to add poller: %v", err)
	}

	return nil
}
