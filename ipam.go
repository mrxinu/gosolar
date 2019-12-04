package gosolar

import (
	"encoding/json"
	"log"
)

// IPAddress is a struct for basic unmarshalling of an IPAddress object
type IPAddress struct {
	Address string `json:"DisplayName"`
	Status  string `json:"Status"`
	//TransientCount string `json"Transient"`
}

func (c *Client) getIP(subnetName string) IPAddress {

	query := `
		SELECT TOP 1
			I.Status
			,I.DisplayName
		FROM IPAM.IPNode I
		WHERE Status = 2
		AND I.Subnet.DisplayName = @subnetName
	`

	// build a map that will hold the parameters for the query above
	parameters := map[string]interface{}{
		"subnetName": subnetName,
	}

	// run the query without with the parameters map above
	res, err := c.Query(query, parameters)

	/* Another possible query is this
	SELECT R.Address as SubnetAddress,
		   R.CIDR,
		   R.FriendlyName,
		   R.PercentUsed,
		(SELECT TOP 1
			I2.IpAddress
		FROM IPAM.IPNode as I2
		WHERE I2.Status=2
		AND I2.SubnetId = R.GroupID ) AS FreeIpAddress
	FROM IPAM.GroupReport as R
	WHERE R.GroupType='8' // This group is "subnets"
	*/

	if err != nil {
		log.Fatal(err)
	}

	var ip IPAddress

	// This should catch an empty ip.
	if err := json.Unmarshal(res, &ip); err != nil {
		log.Fatal(err)
	}

	return ip
}

// Set-SwisObject $swis -Uri 'swis://localhost/Orion/IPAM.IPNode/IPAddress=1.1.1.1' -Properties @{ Alias = 'test1' }
//Invoke-SwisVerb $swis IPAM.SubnetManagement GetFirstAvailableIp @("199.10.1.0", "24")
//Invoke-SwisVerb $swis IPAM.SubnetManagement ChangeIPStatus  @("199.10.1.1", "Used")
//Update: Set-SwisObject $swis -Uri 'swis://localhost/Orion/IPAM.IPNode/IpNodeId=2' -Properties @{ Status = 'Used', Comment = "Reserved by terraform." }
func (c *Client) reserveIP(ipAddress string) {

}
