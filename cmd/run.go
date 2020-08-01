/*
Copyright © 2020 NAME HERE venkyvb@gmail.com

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
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/venkyvb/msft-xl-graph-test/internal"
)



// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the MSFT Graph Excel API tests",
	Long: `Run the MSFT Graph Excel API tests`,
	Run: func(cmd *cobra.Command, args []string) {

		var config apiutils.Config
		err := viper.Unmarshal(&config)
		if err != nil {
			log.Fatalf("unable to decode into struct, %v", err)
		}

		if len(config.InputParams) == 0 {
			config.InputParams = apiutils.GetDefaultInput()
		}

		apiutils.RunTests(config)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
