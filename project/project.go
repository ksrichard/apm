package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ksrichard/apm/util"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var ProjectDetailsFileName string = "apm.json"

func GetProjectDir(cmd *cobra.Command) (string, error) {
	return cmd.Flags().GetString("project-dir")
}

func GetProjectDetails(cmd *cobra.Command) (*ProjectDetails, error) {
	var result ProjectDetails
	projectDir, err := GetProjectDir(cmd)
	if err != nil {
		return nil, err
	}
	jsonFilePath := fmt.Sprintf("%s/%s", projectDir, ProjectDetailsFileName)
	if util.FileExists(jsonFilePath) {
		jsonFile, err := ioutil.ReadFile(jsonFilePath)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(jsonFile, &result)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New(fmt.Sprintf("'%s' not found!", jsonFilePath))
	}
	return &result, nil
}

func UpdateProjectDetails(cmd *cobra.Command, details *ProjectDetails) error {
	projectDir, err := GetProjectDir(cmd)
	if err != nil {
		return err
	}
	jsonFilePath := fmt.Sprintf("%s/%s", projectDir, ProjectDetailsFileName)
	if util.FileExists(jsonFilePath) {
		fileData, err := json.MarshalIndent(details, "", "    ")
		if err != nil {
			return err
		}
		err = os.Remove(jsonFilePath)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(jsonFilePath, fileData, os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		return errors.New(fmt.Sprintf("'%s' not found!", jsonFilePath))
	}
	return nil
}
