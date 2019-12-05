package gosolar

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// IPAddress is a struct for basic unmarshalling of an IPAddress object
type IPAddress struct {
	IPNodeID     int    `json:"IpNodeId"`
	Address      string `json:"IPAddress"`
	Status       int    `json:"Status"`
	StatusString string `json:StatusString`
	Comments     string `json:Comments`
	//TransientCount string `json"Transient"`
}

var statuses = []string{"Used", "Available", "Reserved", "Transient", "Blocked"}

// GetFirstAvailableIP returns the first available ip in the subnet
func (c *Client) GetFirstAvailableIP(subnetAddress string, subnetCIDR string) IPAddress {

	// We need to format the body as an array...
	body := []string{
		subnetAddress,
		subnetCIDR,
	}
	res, err := c.Invoke("IPAM.SubnetManagement", "GetFirstAvailableIp", body)

	// We have to do some fancy slicing to get rid of quotes. Fuck this api.
	bodyString := string(res)[1 : len(string(res))-1]

	if err != nil {
		log.Infof("ResponseString %s", bodyString)
		log.Fatal(err)
	}

	var ip IPAddress

	ip = IPAddress{
		Address: bodyString,
	}

	return ip
}

// GetIP returns a full ip address object for ips
func (c *Client) GetIP(ipAddress string) IPAddress {
	query := `SELECT TOP 1 
				IpNodeId, 
				IPAddress, 
				Status, 
				Comments 
			FROM IPAM.IpNode WHERE IPAddress = @ipAddress`
	parameters := map[string]interface{}{
		"ipAddress": ipAddress,
	}
	res, err := c.Query(query, parameters)
	// run the query without with the parameters map above
	bodyString := string(res)

	if err != nil {
		log.Infof("ResponseString %s", bodyString)
		log.Fatal(err)
	}

	var ip []IPAddress

	if err := json.Unmarshal(res, &ip); err != nil {
		log.Infof("ResponseString %s", bodyString)
		log.Fatal(err)
	}

	// return empty if not found.
	if len(ip) < 1 {
		return IPAddress{}
	}
	ip[0].StatusString = statuses[ip[0].Status-1]
	return ip[0]
}

// Set-SwisObject $swis -Uri 'swis://localhost/Orion/IPAM.IPNode/IPAddress=1.1.1.1' -Properties @{ Alias = 'test1' }
//Invoke-SwisVerb $swis IPAM.SubnetManagement GetFirstAvailableIp @("199.10.1.0", "24")
//Invoke-SwisVerb $swis IPAM.SubnetManagement ChangeIPStatus  @("199.10.1.1", "Used")
//Update: Set-SwisObject $swis -Uri 'swis://localhost/Orion/IPAM.IPNode/IpNodeId=2' -Properties @{ Status = 'Used', Comment = "Reserved by terraform." }

// ReserveIP will set the IP Status to "Used"
func (c *Client) ReserveIP(ipAddress string) IPAddress {
	// We need to format the body as an array...
	body := []string{
		ipAddress,
		"Used",
	}
	result, err := c.Invoke("IPAM.SubnetManagement", "ChangeIPStatus", body)

	resultString := string(result)
	if err != nil {
		log.Info(resultString)
		log.Fatal(err)
	}
	return c.GetIP(ipAddress)
}

// ReleaseIP will set the IP Status to "Unused"
func (c *Client) ReleaseIP(ipAddress string) IPAddress {
	// We need to format the body as an array...
	body := []string{
		ipAddress,
		"Available",
	}
	result, err := c.Invoke("IPAM.SubnetManagement", "ChangeIPStatus", body)
	resultString := string(result)
	if err != nil {
		log.Info(resultString)
		log.Fatal(err)
	}
	return c.GetIP(ipAddress)
}

// CommentOnIPNode puts comments on a ip node object
// https://localhost:17778/SolarWinds/InformationService/v3/Json/swis://--SERVERNAME--/Orion/IPAM.IPNode/IPNodeID=%s
func (c *Client) CommentOnIPNode(ipAddress string, comment string) IPAddress {
	ipAddr := c.GetIP(ipAddress)
	body := map[string]interface{}{
		"Comments": comment,
	}
	log.Info(ipAddr)
	uri := fmt.Sprintf("swis://localhost/Orion/IPAM.IPNode/IpNodeId=%d", ipAddr.IPNodeID)
	log.Info(uri)
	result, err := c.Update(uri, body)
	resultString := string(result)
	if err != nil {
		log.Info(resultString)
		log.Fatal(err)
	}
	return c.GetIP(ipAddress)
}
