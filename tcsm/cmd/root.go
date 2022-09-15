/*
Copyright © 2022 Arush Salil

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
	"context"
	"os"

	"github.com/edjx/opsutil/tcsm/pkg"
	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

var token, workspaceID, address, stateVersion string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tcsm",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		config := &tfe.Config{
			Address: address,
			Token:   token,
		}

		client, err := tfe.NewClient(config)
		if err != nil {
			panic(err)
		}

		err = pkg.RollbackToSpecificVersion(stateVersion, ctx, client, workspaceID)
		if err != nil {
			panic(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&token, "token", "t", "", "Terraform Token")
	rootCmd.MarkFlagRequired("token")
	rootCmd.Flags().StringVarP(&workspaceID, "workspace-id", "w", "", "Workspace ID")
	rootCmd.MarkFlagRequired("workspace-id")
	rootCmd.Flags().StringVarP(&address, "address", "a", "https://app.terraform.io", "Address of TFE host, defaults to TFC.")
	rootCmd.Flags().StringVarP(&stateVersion, "state-version", "s", "", "Version of state to rollback to")
	rootCmd.MarkFlagRequired("state-version")
}
