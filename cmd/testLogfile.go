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
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/unikumAB/logmerge/formats"
)

// testLogfileCmd represents the testLogfile command
var testLogfileCmd = &cobra.Command{
	Use:   "testLogfile",
	Short: "Parses a logfile and outputs some statistics and not parsable log lines",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("testLogfile called on " + args[0])
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Printf("Failed to file file %v", args[0])
			os.Exit(1)
		}
		defer checkedClose(file)
		parser, err := formats.NewVCombinedParser()
		if err != nil {
			panic("Failed to create Parser")
		}
		scanner := bufio.NewScanner(file)
		total, good, bad := 0, 0, 0
		for scanner.Scan() {
			total++
			text := scanner.Text()
			_, err := parser.ParseLine(text)
			if err != nil {
				bad++
				if bad < 10 {
					fmt.Printf("Bad: %v\n", text)
				}
				continue
			}
			good++
		}
		fmt.Printf("Good:\t%v\nBad:\t%v\nTotal:\t%v\n", good, bad, total)
	},
}

func checkedClose(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		log.Fatalf("Failed to close: %v", err)
	}
}
func init() {
	rootCmd.AddCommand(testLogfileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testLogfileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testLogfileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
