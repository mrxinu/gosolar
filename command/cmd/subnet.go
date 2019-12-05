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
)

var vlan string

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
		name, _ := cmd.Flags().GetString("name")

		client := GetClient(cmd, args)
		if 
		resultSubnet, _ := json.Marshal(client.GetSubnet(name))
		fmt.Println(string(resultSubnet))
	},
}

func init() {
	rootCmd.AddCommand(subnetCmd)
	subnetCmd.AddCommand(findSubnet)
	subnetCmd.AddCommand(listSubnets)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// subnetCmd.PersistentFlags().String("foo", "", "A help for foo")
	findSubnet.Flags().StringP("name", "n", "", "Subnet name")

	findSubnet.Flags().StringVarP(&vlan, "vlan", "", "", "Subnet vlan")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// subnetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
