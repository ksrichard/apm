package service

import (
	"errors"
	"fmt"
	"github.com/ksrichard/apm/arduino"
	"github.com/ksrichard/apm/util"
	"strings"
)

func SelectLibrary(cli *arduino.ArduinoCli) (string, string, error) {
	libName := ""
	libVersion := ""

	selectedLibName, err := util.AutoCompleteInput("Library search", "Search again...", "Cancel",
		func(query string) (m map[string]interface{}, err error) {
			result := make(map[string]interface{})
			libs, err := cli.SearchLibrary(query)
			if err != nil {
				return nil, err
			}
			for _, lib := range libs {
				libTitle := fmt.Sprintf("%s (%s) - %s", lib.Name, lib.Latest.Version, lib.Latest.Author)
				result[libTitle] = lib.Name
			}
			return result, nil
		})
	if err != nil {
		return "", "", err
	}
	libName = selectedLibName.(string)

	// select library version
	libVersionOptions := make(map[string]interface{})
	libs, err := cli.SearchLibrary(libName)
	if err != nil {
		return "", "", err
	}
	for _, lib := range libs {
		if lib.Name == libName {
			for _, release := range lib.Releases {
				libVersionOptions[release.Version] = release.Version
			}
			break
		}
	}
	selectedLibVersion, err := util.Select("Select library version", []string{"latest"}, libVersionOptions)
	if err != nil {
		return "", "", err
	}
	libVersion = selectedLibVersion.(string)

	return libName, libVersion, nil
}

func CheckIfLibraryValid(cli *arduino.ArduinoCli, libName string, libVersion string, maxHints int) (string, error) {
	libs, err := cli.SearchLibrary(libName)
	if err != nil {
		return "", err
	}

	hints := []string{}
	libAllVersions := []string{"latest"}

	finalLibName := libName
	foundLibName := false
	foundLibVersion := false
	if strings.ToLower(libVersion) == "latest" {
		foundLibVersion = true
	}
	for _, lib := range libs {
		if strings.ToLower(lib.Name) == strings.ToLower(libName) {
			foundLibName = true
			finalLibName = lib.Name
			for _, release := range lib.Releases {
				libAllVersions = append(libAllVersions, release.Version)
				if release.Version == libVersion {
					foundLibVersion = true
				}
			}
		}

		// add hint if possible
		if (strings.Contains(strings.ToLower(lib.Name), strings.ToLower(libName)) ||
			strings.Contains(lib.Latest.Sentence, libName)) && len(hints) <= maxHints {
			hints = append(hints, fmt.Sprintf("%s@%s", lib.Name, lib.Latest.Version))
		}
	}

	if !foundLibName {
		fmt.Printf("Unknown library name '%s'!\n", libName)
		if len(hints) > 0 {
			fmt.Printf("You can may check also:\n %s \n", strings.Join(hints, "\n"))
		}
		return "", errors.New(fmt.Sprintf("Unknown library name '%s'!", libName))
	}

	if !foundLibVersion {
		fmt.Printf("Unknown library version '%s'!\n", libVersion)
		fmt.Printf("You can use the following versions:\n %s", strings.Join(libAllVersions, "\n"))
		return "", errors.New(fmt.Sprintf("Unknown library version '%s'!", libVersion))
	}

	return finalLibName, nil
}


