// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stobias123/gosolar"
	log "github.com/sirupsen/logrus"
)

var vlanName string

// subnetCmd represents the subnet command
var subnetCmd = &cobra.Command{
	Use:   "subnet",
	Short: "A brief description of your command",
	Long:  `Commands related to subnets can be found here.`,
}

// subnetCmd represents the subnet command
var listSubnets = &cobra.Command{
	Use:   "list subnet <SubnetName>",
	Short: "List all subnets in orion",
	Long:  `List all subnets in orion.`,
	Run: func(cmd *cobra.Command, args []string) {

		client := GetClient(cmd, args)
		resultSubnet, _ := json.Marshal(client.ListSubnets())
		fmt.Println(string(resultSubnet))
	},
}

// findSubnet is a sub command that searches for a subnet by name
var findSubnet = &cobra.Command{
	Use:   "find subnet <SubnetName>",
	Short: "Find a subnet by name",
	Long:  `Find subnet by name.`,
	Run: func(cmd *cobra.Command, args []string) {
		var subnet gosolar.Subnet
		
		client := GetClient(cmd, args)
		subnetName, _ := cmd.Flags().GetString("name")
		
		if len(subnetName) > 1 {
			subnet = client.GetSubnet(subnetName)
		} else if len(vlanName) > 1 {
			subnet = client.GetSubnetByVLAN(vlanName)
		} else {
			log.Errorf("Provide either subnet_name or vlan")
		}

		resultSubnet, _ := json.Marshal(subnet)
		fmt.Println(string(resultSubnet))
	},
}

func init() {
	rootCmd.AddCommand(subnetCmd)
	subnetCmd.AddCommand(findSubnet)
	subnetCmd.AddCommand(listSubnets)
	
	findSubnet.Flags().StringP("name", "n", "", "Subnet name")
	findSubnet.Flags().StringVarP(&vlanName, "vlan", "", "", "Subnet vlan")

}
