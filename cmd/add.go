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
	"github.com/ksrichard/apm/service"
	"github.com/ksrichard/apm/util"
	"github.com/spf13/cobra"
	"strings"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Example: "apm add\napm add OneWire@2.3.5\napm add onewire\napm add onewire@latest",
	Short: "Adding new libraries to the project",
	Long:  `Adding new libraries to the Arduino project`,
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

		// add git repository
		gitRepo, err := cmd.Flags().GetString("git")
		if err != nil {
			return err
		}
		if strings.TrimSpace(gitRepo) != "" {
			return addGitRepoDep(cli, cmd, gitRepo, details)
		}

		// add zip library
		zipFile, err := cmd.Flags().GetString("zip")
		if err != nil {
			return err
		}
		if strings.TrimSpace(zipFile) != "" && util.FileExists(zipFile) {
			return addZipDep(cli, cmd, zipFile, details)
		}
		if strings.TrimSpace(zipFile) != "" && !util.FileExists(zipFile) {
			return errors.New(fmt.Sprintf("'%s' not found!", zipFile))
		}

		// add library
		libName := ""
		libVersion := ""
		if len(args) < 1 { // we do not have any library set
			fmt.Println("No library provided...")
			libName, libVersion, err = service.SelectLibrary(cli)
			if err != nil {
				return err
			}
		} else { // we have library set
			libNameWithVersion := args[0]
			nameAndVer := strings.Split(libNameWithVersion, "@")

			// we have too many @ chars
			if len(nameAndVer) > 2 {
				return errors.New("please provide the library in the following form: LIBRARY_NAME or LIBRARY_NAME@VERSION")
			}

			// only library name provided
			if len(nameAndVer) == 1 {
				libName = nameAndVer[0]
				libVersion = "latest"
			}

			// library name and version provided
			if len(nameAndVer) == 2 {
				libName = nameAndVer[0]
				libVersion = nameAndVer[1]
			}

			// check validity and set library name to original one
			libName, err = service.CheckIfLibraryValid(cli, libName, libVersion, 5)
			if err != nil {
				return err
			}
		}

		// update changes
		fmt.Printf("Adding %s@%s...\n", libName, libVersion)
		hasDep := false
		for i, dep := range details.Dependencies {
			if dep.Library == libName {
				hasDep = true
				details.Dependencies[i].Version = libVersion
			}
		}
		if !hasDep {
			details.Dependencies = append(details.Dependencies, project.ProjectDependency{
				Library: libName,
				Version: libVersion,
			})
		}

		// check if we have any dependency mismatch with current libs
		for _, dep := range details.Dependencies {
			err = cli.CheckDependencyVersionMismatch(dep, details)
			if err != nil {
				return err
			}
		}

		// update project file
		err = project.UpdateProjectDetails(cmd, details)
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
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringP("git", "g", "", "Library from Git repository")
	addCmd.Flags().StringP("zip", "z", "", "Library from ZIP file")
}

func addGitRepoDep(cli *arduino.ArduinoCli, cmd *cobra.Command, gitRepo string, details *project.ProjectDetails) error {
	fmt.Printf("Adding %s...\n", gitRepo)
	hasDep := false
	for i, dep := range details.Dependencies {
		if dep.Git == gitRepo {
			hasDep = true
			details.Dependencies[i].Git = gitRepo
		}
	}
	if !hasDep {
		details.Dependencies = append(details.Dependencies, project.ProjectDependency{
			Git: gitRepo,
		})
	}

	// check if we have any dependency mismatch with current libs
	for _, dep := range details.Dependencies {
		err := cli.CheckDependencyVersionMismatch(dep, details)
		if err != nil {
			return err
		}
	}

	// update project file
	err := project.UpdateProjectDetails(cmd, details)
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
}

func addZipDep(cli *arduino.ArduinoCli, cmd *cobra.Command, zipFile string, details *project.ProjectDetails) error {
	fmt.Printf("Adding %s...\n", zipFile)
	hasDep := false
	for i, dep := range details.Dependencies {
		if dep.Zip == zipFile {
			hasDep = true
			details.Dependencies[i].Zip = zipFile
		}
	}
	if !hasDep {
		details.Dependencies = append(details.Dependencies, project.ProjectDependency{
			Zip: zipFile,
		})
	}

	// check if we have any dependency mismatch with current libs
	for _, dep := range details.Dependencies {
		err := cli.CheckDependencyVersionMismatch(dep, details)
		if err != nil {
			return err
		}
	}

	// update project file
	err := project.UpdateProjectDetails(cmd, details)
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
}
