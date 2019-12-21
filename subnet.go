package gosolar

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

// Subnet this is the json representation of subnet
type Subnet struct {
	SubnetID       int    `json:"SubnetId"`
	Address        string `json:"Address"`
	CIDR           int    `json:"CIDR"`
	FriendlyName   string `json:"FriendlyName"`
	DisplayName    string `json:"DisplayName"`
	AvailableCount int    `json:"AvailableCount"`
	ReservedCount  int    `json:"ReservedCount"`
	UsedCount      int    `json:"UsedCount"`
	TotalCount     int    `json:"totalCount"`
	Comments       string `json:"Comments"`
	VLAN           string `json:"VLAN"`
	AddressMask    string `json:"AddressMask"`
}

// GetSubnet Gets a subnet by display name.
func (c *Client) GetSubnet(subnetName string) Subnet {
	query := `SELECT TOP 1 Address, 
						CIDR, 
						AddressMask, 
						DisplayName, 
						FriendlyName, 
						TotalCount, 
						UsedCount, 
						AvailableCount, 
						ReservedCount, 
						TransientCount,
						VLAN,
						StatusName 
					FROM IPAM.Subnet
					WHERE DisplayName = @name`
	
	parameters := map[string]interface{}{
		"name": subnetName,
	}

	res, err := c.Query(query, parameters)

	var subnet []Subnet
	bodyString := string(res)

	if err != nil {
		log.Infof("ResponseString %s", bodyString)
		log.Fatal(err)
	}

	if err := json.Unmarshal(res, &subnet); err != nil {
		log.Infof("ResponseString %s", bodyString)
		log.Fatal(err)
	}
	if len(subnet) < 1 {
		return Subnet{}
	}

	return subnet[0]
}


// GetSubnetByAddress Gets a subnet using the Address as a parameter
func (c *Client) GetSubnetByAddress(vlan string) Subnet {
	query := `SELECT TOP 1 Address, 
						CIDR, 
						AddressMask, 
						DisplayName, 
						FriendlyName, 
						TotalCount, 
						UsedCount, 
						AvailableCount, 
						ReservedCount, 
						TransientCount,
						VLAN,
						StatusName 
					FROM IPAM.Subnet
					WHERE Address = @address`

	parameters := map[string]interface{}{
		"address": Address,
	}

	res, err := c.Query(query, parameters)

	var subnet []Subnet
	bodyString := string(res)

	if err != nil {
		log.Infof("ResponseString %s", bodyString)
		log.Fatal(err)
	}

	if err := json.Unmarshal(res, &subnet); err != nil {
		log.Infof("ResponseString %s", bodyString)
		log.Fatal(err)
	}
	if len(subnet) < 1 {
		return Subnet{}
	}

	return subnet[0]
}

// GetSubnetByVLAN Gets a subnet using the VLAN as a parameter
func (c *Client) GetSubnetByVLAN(vlan string) Subnet {
	query := `SELECT TOP 1 Address, 
						CIDR, 
						AddressMask, 
						DisplayName, 
						FriendlyName, 
						TotalCount, 
						UsedCount, 
						AvailableCount, 
						ReservedCount, 
						TransientCount,
						VLAN,
						StatusName 
					FROM IPAM.Subnet
					WHERE VLAN = @vlan`

	parameters := map[string]interface{}{
		"vlan": vlan,
	}

	res, err := c.Query(query, parameters)

	var subnet []Subnet
	bodyString := string(res)

	if err != nil {
		log.Infof("ResponseString %s", bodyString)
		log.Fatal(err)
	}

	if err := json.Unmarshal(res, &subnet); err != nil {
		log.Infof("ResponseString %s", bodyString)
		log.Fatal(err)
	}
	if len(subnet) < 1 {
		return Subnet{}
	}

	return subnet[0]
}

// ListSubnets Lists subnets.
func (c *Client) ListSubnets() []Subnet {
	query := `SELECT	Address, 
						CIDR, 
						AddressMask, 
						DisplayName, 
						FriendlyName, 
						TotalCount, 
						UsedCount, 
						AvailableCount, 
						ReservedCount, 
						TransientCount, 
							StatusName 
					FROM IPAM.Subnet`

	res, err := c.Query(query, nil)

	var subnets []Subnet
	bodyString := string(res)

	if err != nil {
		log.Info("Couldnt unmarshal responseString %s", bodyString)
		log.Fatal(err)
	}

	if err := json.Unmarshal(res, &subnets); err != nil {
		log.Info("Couldnt unmarshal responseString %s", bodyString)
		log.Fatal(err)
	}
	return subnets
}
