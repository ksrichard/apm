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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ksrichard/apm/project"
	"github.com/ksrichard/apm/util"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init APM project",
	Long: `Init Arduino Package Manager project`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectDir, err := project.GetProjectDir(cmd)
		if err != nil {
			return err
		}
		jsonFilePath := fmt.Sprintf("%s/%s", projectDir, project.ProjectDetailsFileName)
		if util.FileExists(jsonFilePath) {
			return errors.New(fmt.Sprintf("'%s' is already initialized", projectDir))
		} else {
			details := project.ProjectDetails{
				Board:        &project.ProjectBoard{},
				Dependencies: []project.ProjectDependency{},
			}
			fileData, err := json.MarshalIndent(details, "", "    ")
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(jsonFilePath, fileData, os.ModePerm)
			if err != nil {
				return err
			}
			return nil
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
