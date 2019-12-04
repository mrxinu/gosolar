package gosolar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type TestSuite struct {
	suite.Suite
	client *Client
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *TestSuite) SetupTest() {
	hostname := "localhost"
	username := "bleh"
	password := "bleh"
	suite.client = NewClient(hostname, username, password, true)
}

func (suite *TestSuite) TestFetchSubnet() {

	// assert equality
	//assert.Equal(t, 123, 123, "they should be equal")
	subnet := suite.client.GetSubnet("Subnet1")
	nullSubnetList := []Subnet{}
	assert.Equal(suite.T(), nullSubnetList, subnet)
}

func (suite *TestSuite) TestFetchNamedSubnet() {

	// assert equality
	//assert.Equal(t, 123, 123, "they should be equal")
	subnet := suite.client.GetSubnet("test")
	expectedSubnet := Subnet{
		SubnetID:       1234,
		Address:        "10.199.152.0",
		CIDR:           "23",
		FriendlyName:   "test subnet",
		DisplayName:    "test subnet",
		AvailableCount: 200,
		ReservedCount:  2,
		UsedCount:      181,
		TotalCount:     512,
		Comments:       "NFS - VLAN 410",
		VLAN:           410,
		AddressMask:    "255.255.254.0",
	}
	assert.Equal(suite.T(), expectedSubnet, subnet)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
