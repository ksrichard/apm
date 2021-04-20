/*
Copyright © 2021 Richard Klavora <klavorasr@gmail.com>

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
	"github.com/ksrichard/apm/arduino"
	"github.com/ksrichard/apm/project"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install dependencies of project",
	Long: `Install dependencies of the Arduino project`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// project details
		details, err := project.GetProjectDetails(cmd)
		if err != nil {
			return err
		}

		cli := &arduino.ArduinoCli{}
		err = cli.Init()
		if err != nil {
			return err
		}
		defer cli.Destroy()

		// install board core package
		if details.Board != nil && details.Board.Package != "" {
			err = cli.InstallBoardCore(details)
			if err != nil {
				return err
			}
		}

		// install dependencies
		if details.Dependencies != nil && len(details.Dependencies) > 0 {
			err = cli.InstallDependencies(details)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
