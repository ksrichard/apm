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
	"log"
	"strings"

	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use: "remove",
	Example: "apm remove\n" +
		"apm remove OneWire\n" +
		"apm remove onewire\n" +
		"apm remove \"robot control\"\n" +
		"apm remove \"Robot Control\"\n" +
		"apm remove \"https://github.com/jandrassy/ArduinoOTA\"\n" +
		"apm remove https://github.com/jandrassy/ArduinoOTA\n" +
		"apm remove ArduinoOTA.zip\n" +
		"apm remove \"ArduinoOTA.zip\"",
	Short: "Remove library from the project",
	Long:  `Remove library from the Arduino project`,
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

		var libToRemove *project.ProjectDependency = nil
		// no library provided
		if len(args) < 1 {
			fmt.Println("No library provided...")
			items := make(map[string]interface{})
			for _, dep := range details.Dependencies {
				libTitle := ""
				if dep.Library == "" && dep.Git != "" {
					libTitle = dep.Git
				}
				if dep.Library == "" && dep.Zip != "" {
					libTitle = dep.Zip
				}
				if dep.Library != "" && dep.Version != "" {
					libTitle = fmt.Sprintf("%s (%s)", dep.Library, dep.Version)
				}
				items[libTitle] = &dep
			}
			selectedLib, err := util.Select("Select library to remove", []string{"Cancel"}, items)
			if err != nil {
				return err
			}

			// check if cancelled
			switch selectedLib.(type) {
			case string:
				if selectedLib.(string) == "Cancel" {
					return errors.New("cancelled")
				}
				break
			}

			libToRemove = selectedLib.(*project.ProjectDependency)
		} else { // library name provided
			libToRemoveArg := args[0]
			for _, dep := range details.Dependencies {
				if strings.ToLower(dep.Library) == strings.ToLower(libToRemoveArg) ||
					dep.Git == libToRemoveArg ||
					dep.Zip == libToRemoveArg {
					libToRemove = &dep
				}
			}
			if libToRemove == nil {
				return errors.New(fmt.Sprintf("failed to find '%s' library in the project", libToRemoveArg))
			}
		}

		// log removal
		libName := ""
		if libToRemove.Library != "" {
			libName = libToRemove.Library
		}
		if libToRemove.Git != "" {
			libName = libToRemove.Git
		}
		if libToRemove.Zip != "" {
			libName = libToRemove.Zip
		}
		log.Printf("Removing '%s'...", libName)

		// remove from project file
		for i, dep := range details.Dependencies {
			if (dep.Library != "" && dep.Library == libToRemove.Library) ||
				(dep.Git != "" && dep.Git == libToRemove.Git) ||
				(dep.Zip != "" && dep.Zip == libToRemove.Zip) {
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
