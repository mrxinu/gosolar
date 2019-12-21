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
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stobias123/gosolar"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "solarcmd",
	Short: "solarcmd is a CLI to interact with Solarwinds Orion API",
	Long: `solarcmd is a CLI application that provides several common orion functions.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version is a CLI to interact with Solarwinds Orion API",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) { fmt.Println("1.0")},
}

func init() {
	cobra.OnInitialize(initConfig)
	
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.solar.yaml)")
	rootCmd.PersistentFlags().StringP("server", "s", viper.GetString("ORION_SERVER"), "Set Orion Server")
	rootCmd.PersistentFlags().StringP("debug", "", viper.GetString("ORION_DEBUG"), "Set debug level")
	rootCmd.PersistentFlags().StringP("username", "u", viper.GetString("ORION_USERNAME"), "Set Orion Username")
	rootCmd.PersistentFlags().StringP("password", "p", viper.GetString("ORION_PASSWORD"), "Set Orion Password")
	rootCmd.PersistentFlags().BoolP("ssl", "", viper.GetBool("ORION_USE_SSL"), "Use SSL")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".solar" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".solar")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// GetClient is a convenience function to create a client object with provided strings.
func GetClient(cmd *cobra.Command, args []string) *gosolar.Client {
	hostname, _ := cmd.Flags().GetString("server")
	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")
	ssl, _ := cmd.Flags().GetBool("ssl")

	return gosolar.NewClient(hostname, username, password, ssl, true)
}
