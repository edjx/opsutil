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
	"os"

	"github.com/arush-sal/doddl/pkg/executor"
	"github.com/arush-sal/doddl/pkg/getter"
	"github.com/arush-sal/doddl/pkg/print"
	"github.com/spf13/cobra"
)

var token, tag, list string
var json, stop bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "doddl",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var l []getter.Droplet = executor.Run(executor.GetDOClient(token), list, tag)
		if stop {
			l = executor.RunStopped(executor.GetDOClient(token), l)
		}

		if json {
			print.JSONPrinter(l)
		} else {
			print.Printer(l)
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
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "DigitalOcean Personal Access Token")
	rootCmd.PersistentFlags().StringVarP(&list, "list", "l", "", "Provide a listing option out of [all-tagged, no-tag, tag, no-whitelist]\nMake sure to provide the tag when using tag option")
	rootCmd.PersistentFlags().StringVarP(&tag, "tag", "a", "", "Tag you want to list the Droplets by")
	rootCmd.PersistentFlags().BoolVarP(&json, "json", "j", false, "Print the output in JSON")
	rootCmd.PersistentFlags().BoolVarP(&stop, "stopped", "s", false, "List all the stopped instaces")
}
