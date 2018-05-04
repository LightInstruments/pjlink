// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"github.com/LightInstruments/pjlink"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Display status of Projector",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if projectorIp == "" {
			fmt.Println("projectorIp has to be specified.")
			os.Exit(1)
		}

		proj := pjlink.NewProjector(projectorIp, password)

		stat, err := proj.GetPowerStatus()
		if err != nil {
			log.Println(err)
		} else {
			log.Println(stat)
		}


		err = proj.TurnOn()
		if err != nil {
			log.Println(err)
		}

		stat, err = proj.GetPowerStatus()
		if err != nil {
			log.Println(err)
		} else {
			log.Println(stat)
		}

		err = proj.TurnOff()
		if err != nil {
			log.Println(err)
		}

		stat, err = proj.GetPowerStatus()
		if err != nil {
			log.Println(err)
		} else {
			log.Println(stat)
		}

	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
