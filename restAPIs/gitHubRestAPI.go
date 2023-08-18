package restAPIs

import (
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"html/template"
	"httpServer/requestsHandlers"
	"math"
	"net/http"
	"strconv"
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
		Validate: nil, // Validation is not required for optional input
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
	case Organization:
		urlString += "+org:"
	case Owner:
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

	totalRepos := len(githubResponse.Items)
	totalPages := int(math.Ceil(float64(totalRepos) / float64(perPage)))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		currentPage := getCurrentPage(r)
		startIndex, endIndex := getPageIndices(currentPage, perPage, len(githubResponse.Items))
		paginatedResponse := githubResponse.Items[startIndex:endIndex]
		pageNumbers := getPageNumbers(currentPage, totalPages)

		tmpl, err := template.ParseFiles("templates/repositories.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		queryParams := getQueryParams(r)

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

func getCurrentPage(r *http.Request) int {
	query := r.URL.Query()
	pageParam := query.Get("page")
	currentPage := 1
	if pageParam != "" {
		currentPage, _ = strconv.Atoi(pageParam)
	}
	return currentPage
}

func getPageIndices(currentPage int, perPage int, totalItems int) (int, int) {
	startIndex := (currentPage - 1) * perPage
	endIndex := startIndex + perPage
	if endIndex > totalItems {
		endIndex = totalItems
	}
	return startIndex, endIndex
}

func getPageNumbers(currentPage int, totalPages int) []int {
	startPage := int(math.Max(float64(currentPage-2), 1))
	endPage := int(math.Min(float64(currentPage+2), float64(totalPages)))

	pageNumbers := make([]int, endPage-startPage+1)
	for i := startPage; i <= endPage; i++ {
		pageNumbers[i-startPage] = i
	}
	return pageNumbers
}

func getQueryParams(r *http.Request) map[string]string {
	queryParams := make(map[string]string)
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			queryParams[key] = values[0]
		}
	}
	return queryParams
}

func (api *GitHubRestAPI) FetchRepositoriesByType() {
	filter, content, phrase, err := api.GetUserInput()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	url := api.BuildUrl(filter, content, phrase)

	response := api.SendGetRequest(url)

	// Define the desired number of repositories per page
	perPage := 12

	api.DisplayResponse(response, perPage)
}
