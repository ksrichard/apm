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
	"errors"
	"fmt"
	"github.com/ksrichard/apm/arduino"
	"github.com/ksrichard/apm/project"
	"github.com/ksrichard/apm/util"
	"strings"

	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove",
	Example: "apm remove\napm remove OneWire\napm remove onewire\napm remove \"robot control\"\napm remove \"Robot Control\"",
	Short:   "Remove library from the project",
	Long:    `Remove library from the Arduino project`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// project details
		details, err := project.GetProjectDetails(cmd)
		if err != nil {
			return err
		}

		// init cli
		cli := &arduino.ArduinoCli{}
		err = cli.Init()
		if err != nil {
			return err
		}
		defer cli.Destroy()

		libToRemove := ""
		// no library provided
		if len(args) < 1 {
			fmt.Println("No library provided...")
			items := make(map[string]interface{})
			for _, dep := range details.Dependencies {
				// TODO: continue
				//if dep.Library == "" && dep.Git != "" {
				//	items[dep.Git] = dep.Git
				//}
				//if dep.Library == "" && dep.Zip != "" {
				//	items[dep.Zip] = dep.Zip
				//}
				if dep.Library != "" && dep.Version != "" {
					items[fmt.Sprintf("%s (%s)", dep.Library, dep.Version)] = dep.Library
				}
			}
			selectedLib, err := util.Select("Select library to remove", []string{"Cancel"}, items)
			if err != nil {
				return err
			}
			if selectedLib.(string) == "Cancel" {
				return errors.New("cancelled")
			}
			libToRemove = selectedLib.(string)
		} else { // library name provided
			libToRemoveArg := args[0]
			for _, dep := range details.Dependencies {
				if strings.ToLower(dep.Library) == libToRemoveArg {
					libToRemove = dep.Library
				}
			}
			if libToRemove == "" {
				return errors.New(fmt.Sprintf("failed to find '%s' in the project", libToRemoveArg))
			}
		}

		// remove from project file
		for i, dep := range details.Dependencies {
			if dep.Library == libToRemove {
				details.Dependencies = removeFromDeps(details.Dependencies, i)
				break
			}
		}

		// update project file
		err = project.UpdateProjectDetails(cmd, details)
		if err != nil {
			return err
		}

		// uninstall dependency
		err = cli.UninstallDependency(libToRemove)
		if err != nil {
			return err
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
	rootCmd.AddCommand(removeCmd)
}

func removeFromDeps(slice []project.ProjectDependency, i int) []project.ProjectDependency {
	return append(slice[:i], slice[i+1:]...)
}
