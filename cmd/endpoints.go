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
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/pkg/configtx"
	"github.com/spf13/cobra"
)

var (
	OrderingEndpoint string
	OrgName          string
)

// endpointsCmd represents the endpoints command
var endpointsCmd = &cobra.Command{
	Use:   "endpoints",
	Short: "Updates orderer endpoint",
	Long: `Adds an orderer's endpoint to an existing channel config transaction. If
	the same endpoint already exist in current configuration, this will be a no-op.
  For example:

configtx orderer update endpoints --orgName Org1 --endpoint 127.0.0.1:8080
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("endpoints called")
		updateEndpoints()
	},
}

func init() {
	updateCmd.AddCommand(endpointsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// endpointsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// endpointsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	localFlags := endpointsCmd.Flags()
	localFlags.StringVarP(&OrderingEndpoint, "endpoint", "", "", "Ordering service endpoint")
	localFlags.StringVarP(&OrgName, "orgName", "", "", "Organization name")
}

func updateEndpoints() {

	config, err := readBlock(ConfigFilePath)
	if err != nil {
		panic(err)
	}

	if len(strings.Split(OrderingEndpoint, ":")) != 2 {
		panic("ordering service endpoint %s is not valid or missing")
	}

	endpoint := strings.Split(OrderingEndpoint, ":")

	port, err := strconv.Atoi(endpoint[1])
	if err != nil {
		panic(err)
	}

	address := configtx.Address{
		Host: endpoint[0],
		Port: port,
	}

	err = config.SetOrdererEndpoint(OrgName, address)
	if err != nil {
		panic(err)
	}
}
