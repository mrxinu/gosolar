package gosolar

import (
	"encoding/json"
	"fmt"
	"testing"
)

// TestRemoveNCMNodes will test the getRemoveNCMNodesRequest function for stability
func TestRemoveNCMNodes(t *testing.T) {
	testGuids := []string{"guid1", "guid2"}
	expectedReq := [][]string{testGuids}
	expectedEndpoint := "Invoke/Cirrus.Nodes/RemoveNodes"
	actualReq, actualEndpoint := getRemoveNCMNodesRequest(testGuids)
	if actualEndpoint != expectedEndpoint {
		t.Fatalf("Invalid endpoint. Expected [%s] recieved [%s].", expectedEndpoint, actualEndpoint)
	}
	err := lazyCompare(expectedReq, actualReq)
	if err != nil {
		t.Fatalf("Invalid request: %v", err)
	}
}

// lazyCompare will take two interfaces and compare them for length and that their marshaled
// forms are equal on the byte level.
func lazyCompare(expected, actual interface{}) error {
	expectedBytes, err := json.Marshal(expected)
	if err != nil {
		return err
	}
	actualBytes, err := json.Marshal(actual)
	if err != nil {
		return err
	}
	elen := len(expectedBytes)
	alen := len(actualBytes)
	if elen != alen {
		return fmt.Errorf("Length discrepency. Expected length [%d] recieved length [%d].", elen, alen)
	}
	for i, actualByte := range actualBytes {
		expectedByte := expectedBytes[i]
		if actualByte != expectedByte {
			return fmt.Errorf("Byte discrepency. Expected [%X] recieved length [%X].", expectedByte, actualByte)
		}
	}
	return nil
}
