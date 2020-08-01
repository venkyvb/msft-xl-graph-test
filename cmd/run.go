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

		accessToken := viper.GetString("accessToken")
		workbookItemID := viper.GetString("workbookItemID")
		noOfIterations := viper.GetInt("noOfIterations")

		apiutils.RunTests(accessToken, workbookItemID, noOfIterations)
	},
}

func init() {
	var accessToken string
	var workbookItemID string
	var noOfIterations int

	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&accessToken, "accessToken", "a", "", "Enter the accessToken")
	runCmd.Flags().StringVarP(&workbookItemID, "workbookItemID", "w", "", "Enter the workbookItemID")
	runCmd.Flags().IntVarP(&noOfIterations, "noOfIterations", "n", 1, "Enter the noOfIterations, default would be 1")	

	viper.BindPFlags(runCmd.Flags())
}
