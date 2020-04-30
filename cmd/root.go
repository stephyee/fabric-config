/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/pkg/configtx"
	"github.com/hyperledger/fabric/protoutil"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile        string
	ConfigFilePath string
	Policy         string
	ACLs           string
	Capability     string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "configtx",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	flags := rootCmd.PersistentFlags()
	flags.StringVarP(&ConfigFilePath, "file", "", "", "A channel configuration in protobuf format")
	flags.StringVarP(&Policy, "policy", "p", "", "Defines rules for access to channels, chaincodes, etc.")
	flags.StringVarP(&ACLs, "acls", "a", "", "Defines access to resources by associating a policy with a resource")
	flags.StringVarP(&Capability, "capability", "", "", "Defines capabilities of a fabric network")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

		// Search config in home directory with name ".configtx" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".configtx")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func readBlock(cfgPath string) (configtx.ConfigTx, error) {
	data, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return configtx.ConfigTx{}, fmt.Errorf("Could not read block %s", cfgPath)
	}

	blk, err := protoutil.UnmarshalBlock(data)
	if err != nil {
		panic(err)
	}

	payload, err := protoutil.UnmarshalPayload(blk.Data.Data[0])
	if err != nil {
		panic(err)
	}

	config := &cb.Config{}
	err = proto.Unmarshal(payload.Data, config)
	if err != nil {
		panic(err)
	}

	return configtx.New(config), nil

}
