package restAPIs

import (
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"html/template"
	"httpServer/utils"
	"math"
	"net/http"
	"time"
)

const (
	GitHubBaseUrl         = "https://api.github.com/"
	DefaultResultsPerPage = 12
	FullTimeLayout        = "2006-01-02T15:04:05Z"
	DisplayTimeLayout     = "2006-01-02 15:04"
)

// GitHubRestAPI implements the RestAPI interface for GitHub repositories
type GitHubRestAPI struct{}

type githubRepoJson struct {
	Items []struct {
		Name  string `json:"name"`
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
		URL          string `json:"html_url"`
		CreationTime string `json:"created_at"`
		Stars        int    `json:"stargazers_count"`
	} `json:"items"`
}

func (api *GitHubRestAPI) GetUserInput() (filter string, content string, phrase string, err error) {
	var input struct {
		Filter  string
		Content string
		Phrase  string
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
		return "", "", "", err
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
		return "", "", "", err
	}

	phraseQuestion := &survey.Question{
		Name:     "Phrase",
		Prompt:   &survey.Input{Message: "Please provide a phrase (optional):"},
		Validate: nil,
	}

	err = survey.Ask([]*survey.Question{phraseQuestion}, &input)
	if err != nil {
		return "", "", "", err
	}

	return input.Filter, input.Content, input.Phrase, nil
}

func (api *GitHubRestAPI) BuildUrl(filter string, content string, phrase string) string {
	var urlString = GitHubBaseUrl + "search/repositories?q="
	if phrase != "" {
		urlString += phrase + "+in:name"
	}
	switch filter {
	case "Organization":
		urlString += "+org:"
	case "Owner":
		urlString += "+user:"
	default:
		urlString += "+org:"
	}
	url := urlString + content
	return url
}

func (api *GitHubRestAPI) DisplayResponse(response *http.Response, perPage int) {
	var githubResponse githubRepoJson

	err := json.NewDecoder(response.Body).Decode(&githubResponse)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for i, repo := range githubResponse.Items {
		t, _ := time.Parse(FullTimeLayout, repo.CreationTime)
		githubResponse.Items[i].CreationTime = t.Format(DisplayTimeLayout)
	}

	totalRepos := len(githubResponse.Items)
	totalPages := int(math.Ceil(float64(totalRepos) / float64(perPage)))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		currentPage := paginationHandler.GetCurrentPage(r)
		startIndex, endIndex := paginationHandler.GetPageIndices(currentPage, perPage, len(githubResponse.Items))
		paginatedResponse := githubResponse.Items[startIndex:endIndex]
		pageNumbers := paginationHandler.GetPageNumbers(currentPage, totalPages)

		tmpl, err := template.ParseFiles("templates/repositories.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		queryParams := paginationHandler.GetQueryParams(r)

		data := struct {
			Repositories []struct {
				Name  string `json:"name"`
				Owner struct {
					Login string `json:"login"`
				} `json:"owner"`
				URL          string `json:"html_url"`
				CreationTime string `json:"created_at"`
				Stars        int    `json:"stargazers_count"`
			}
			CurrentPage int
			PageNumbers []int
			BaseUrl     string
			QueryParams map[string]string
		}{
			Repositories: paginatedResponse,
			CurrentPage:  currentPage,
			PageNumbers:  pageNumbers,
			BaseUrl:      response.Request.URL.Path,
			QueryParams:  queryParams,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Server is running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func (api *GitHubRestAPI) FetchRepositoriesByType() {
	filter, content, phrase, err := api.GetUserInput()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	url := api.BuildUrl(filter, content, phrase)

	response, err := http.Get(url)
	if err != nil {
		println(err)
	}

	api.DisplayResponse(response, DefaultResultsPerPage)
}
