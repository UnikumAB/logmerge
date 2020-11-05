/*
Copyright © 2020 Unikum AB <info@unikum.net>

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
	"github.com/UnikumAB/logmerge/merge"
	"github.com/spf13/cobra"
)

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merges the log files",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		merge.Merge(outputFile, args)
	},
}
var (
	outputFile string
	files      []string
)

func init() {
	rootCmd.AddCommand(mergeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	mergeCmd.Flags().StringArrayVarP(&files, "file", "f", []string{"access.log"}, "Access files to merge")
	mergeCmd.Flags().StringVarP(&outputFile, "output", "o", "access-out.log", "Access files to write to")
	// mergeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mergeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
