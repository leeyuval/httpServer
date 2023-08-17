package main

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"httpServer/restAPIs"
)

func main() {
	var input struct {
		Filter  string
		Content string
	}

	questions := []*survey.Question{
		{
			Name: "Filter",
			Prompt: &survey.Select{
				Message: "Please select a filter:",
				Options: []string{"Organization", "Owner"},
			},
			Validate: survey.Required,
		},
		{
			Name:     "Content",
			Prompt:   &survey.Input{Message: "Please provide an appropriate content:"},
			Validate: survey.Required,
		},
	}

	err := survey.Ask(questions, &input)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var gitHubAPI restAPIs.GitHubRestAPI
	gitHubAPI.FetchRepositoriesByFilter(input.Filter, input.Content)
}
