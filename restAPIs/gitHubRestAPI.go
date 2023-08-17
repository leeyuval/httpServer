package restAPIs

import (
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"html/template"
	"httpServer/requestsHandlers"
	"net/http"
	"time"
)

const Organization = "Organization"
const Owner = "Owner"

const GitHubBaseUrl = "https://api.github.com/"
const organizationFilterExt = "orgs/%s/repos"
const ownerFilterExt = "users/%s/repos"

// GitHubRestAPI implements the RestAPI interface for GitHub repositories
type GitHubRestAPI struct {
	requestsHandlers.BasicRequestsHandler
}

type githubRepoJson struct {
	Name  string `json:"name"`
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
	URL          string    `json:"html_url"`
	CreationTime time.Time `json:"created_at"`
	Stars        int       `json:"stargazers_count"`
}

func (api *GitHubRestAPI) GetUserInput() (filter string, content string, err error) {
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

	err = survey.Ask(filterQuestion, &input)
	if err != nil {
		return "", "", err
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
		return "", "", err
	}

	return input.Filter, input.Content, nil
}

func (api *GitHubRestAPI) BuildUrl(filter string, content string) string {
	var urlString = GitHubBaseUrl
	switch filter {
	case Organization:
		urlString += organizationFilterExt
	case Owner:
		urlString += ownerFilterExt
	default:
		urlString += organizationFilterExt
	}
	url := fmt.Sprintf(urlString, content)
	return url
}

func (api *GitHubRestAPI) DisplayResponse(response *http.Response) {
	var githubResponse []githubRepoJson

	err := json.NewDecoder(response.Body).Decode(&githubResponse)
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

		err = tmpl.Execute(w, githubResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Server is running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func (api *GitHubRestAPI) FetchRepositoriesByFilter() {

	filter, content, err := api.GetUserInput()
	if err != nil {
		fmt.Println("Error:", err)
	}

	url := api.BuildUrl(filter, content)

	response := api.SendGetRequest(url)

	api.DisplayResponse(response)

}
