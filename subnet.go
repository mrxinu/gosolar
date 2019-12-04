package gosolar

import (
	"encoding/json"
	"log"
)

type Subnet struct {
	Address        string `json:"Address"`
	CIDR           string `json:"CIDR"`
	Comments       string `json:"Comments"`
	AddressMask    string `json:"AddressMask"`
	DisplayName    string `json:"DisplayName"`
	FriendlyName   string `json:"FriendlyName"`
	TotalCount     int    `json:"totalCount"`
	UsedCount      int    `json:"UsedCount"`
	AvailableCount int    `json:"AvailableCount"`
	ReservedCount  int    `json:"ReservedCount"`
	//TransientCount string `json"Transient"`
}

func (c *Client) GetSubnet(subnetName string) []Subnet {
	query := `SELECT	Address, 
						CIDR, 
						AddressMask, 
						DisplayName, 
						FriendlyName, 
						Reserved, 
						TotalCount, 
						UsedCount, 
						AvailableCount, 
						ReservedCount, 
						TransientCount, 
							StatusName 
					FROM IPAM.Subnets 
					WHERE DisplayName == @name`

	parameters := map[string]interface{}{
		"name": subnetName,
	}

	res, err := c.QueryRow(query, parameters)

	if err != nil {
		log.Fatal(err)
	}

	var subnets []Subnet

	if err := json.Unmarshal(res, &subnets); err != nil {
		log.Fatal(err)
	}

	if len(subnets) < 1 {
		log.Print("[INFO] No subnets matching that name found.")
	}

	return subnets
}
