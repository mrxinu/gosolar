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
	"fmt"

	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cidr string
var address string
var comment string

// ipCmd represents the ip command
var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "Commands related to IPAM.",
	Long: `The ip command command handles common interactions with the IPAM API
Usage:`,
}

var ipGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get an IP Address from subnet.",
	Long: `Finds an IP from a given subnet address. 
  For example:

  solarcmd ip get --subnet_address 10.200.20.0 --cidr 20`,
	Run: func(cmd *cobra.Command, args []string) {
		client := GetClient(cmd, args)
		subnetAddress, err := cmd.Flags().GetString("subnet_address")
		if err != nil {
			log.Fatal(err)
		}
		result := client.GetFirstAvailableIP(subnetAddress, cidr)

		resultIP, _ := json.Marshal(result)
		fmt.Println(string(resultIP))
	},
}

var reserveCmd = &cobra.Command{
	Use:   "reserve",
	Short: "Reserve an IP Address from subnet.",
	Long: `Marks the given IP as reserveCmd.
  For example:

  solarcmd ip reserve --address 192.168.4.5`,
	Run: func(cmd *cobra.Command, args []string) {
		client := GetClient(cmd, args)

		result := client.ReserveIP(address)
		resultIP, _ := json.Marshal(result)
		fmt.Println(string(resultIP))
	},
}

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Release an IP Address from subnet.",
	Long: `Marks the given IP as released.
  For example:

  solarcmd ip release --address 192.168.4.5`,
	Run: func(cmd *cobra.Command, args []string) {
		client := GetClient(cmd, args)

		result := client.ReleaseIP(address)
		resultIP, _ := json.Marshal(result)
		fmt.Println(string(resultIP))
	},
}

var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Comment an IP Address from subnet.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

  solarcmd ip comment --subnet 192.168.4.5 --comment "hello from solarcmd!"`,
	Run: func(cmd *cobra.Command, args []string) {
		client := GetClient(cmd, args)

		result := client.CommentOnIPNode(address, comment)
		resultIP, _ := json.Marshal(result)
		fmt.Println(string(resultIP))
	},
}

var lookupCmd = &cobra.Command{
	Use:   "lookup",
	Short: "Adds comment to IP.Node",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

  solarcmd ip comment -a 192.168.4.5 -comment "hello"`,
	Run: func(cmd *cobra.Command, args []string) {
		client := GetClient(cmd, args)

		result := client.GetIP(address)
		resultIP, _ := json.Marshal(result)
		fmt.Println(string(resultIP))
	},
}

func init() {
	rootCmd.AddCommand(ipCmd)
	ipCmd.AddCommand(ipGetCmd)
	ipCmd.AddCommand(lookupCmd)
	ipCmd.AddCommand(reserveCmd)
	ipCmd.AddCommand(releaseCmd)
	ipCmd.AddCommand(commentCmd)

	ipGetCmd.Flags().StringP("subnet_address", "n", "", "Subnet address")
	ipGetCmd.Flags().StringVarP(&cidr, "cidr", "c", "24", "CIDR Mask")

	reserveCmd.Flags().StringVarP(&address, "address", "a", "", "Address")

	lookupCmd.Flags().StringVarP(&address, "address", "a", "", "Address")

	commentCmd.Flags().StringVarP(&address, "address", "a", "", "Address")
	commentCmd.Flags().StringVarP(&comment, "comment", "", "", "Comment")

	releaseCmd.Flags().StringVarP(&address, "address", "a", "", "Address")
}
