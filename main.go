package main

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"html/template"
	"httpServer/restAPIs"
	"net/http"
)

func main() {
	var input struct {
		Filter  string
		Content string
	}

	filterQuestion := []*survey.Question{
		{
			Name: "Filter",
			Prompt: &survey.Select{
				Message: "Please select a filter:",
				Options: []string{"Organization", "Owner"},
			},
			Validate: survey.Required,
		},
	}

	err := survey.Ask(filterQuestion, &input)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	contentPrompts := map[string]string{
		"Organization": "Please provide organization name:",
		"Owner":        "Please provide owner name:",
	}

	contentQuestion := &survey.Question{
		Name:     "Content",
		Prompt:   &survey.Input{Message: contentPrompts[input.Filter]},
		Validate: survey.Required,
	}

	err = survey.Ask([]*survey.Question{contentQuestion}, &input)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var gitHubAPI restAPIs.GitHubRestAPI
	repositories, err := gitHubAPI.FetchRepositoriesByFilter(input.Filter, input.Content)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/repositories.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, repositories)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Server is running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
