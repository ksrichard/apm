package util

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"strings"
)

func Select(label string, initialItems []string, items map[string]interface{}) (interface{}, error) {
	var promptItems = initialItems
	for k, _ := range items {
		promptItems = append(promptItems, k)
	}

	prompt := promptui.Select{
		Label: label,
		Items: promptItems,
	}

	_, promptResult, err := prompt.Run()

	if err != nil {
		return "", err
	}

	var result interface{}
	for k, v := range items {
		if k == promptResult {
			result = v
		}
	}

	for _, item := range initialItems {
		if item == promptResult {
			result = item
		}
	}

	return result, nil
}

func AutoCompleteInput(title string, searchAgainStr string, cancelStr string, results func(query string) (map[string]interface{}, error)) (interface{}, error) {
	var selectedOption interface{}
	validate := func(input string) error {
		if strings.TrimSpace(input) == "" {
			return errors.New("Input required!")
		}
		return nil
	}
	for {
		// search
		prompt := promptui.Prompt{
			Label:    title,
			Validate: validate,
		}

		query, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// select from result list
		items := make(map[string]interface{})
		res, err := results(query)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if res != nil && len(res) == 0 {
			fmt.Printf("No library found for search query '%s'!\n", query)
			continue
		}
		for k, v := range res {
			items[k] = v
		}

		result, err := Select(title, []string{cancelStr, searchAgainStr}, items)

		if result == cancelStr {
			return "", errors.New("cancelled")
		}

		if result != searchAgainStr {
			selectedOption = result
			break
		}
	}

	return selectedOption, nil
}