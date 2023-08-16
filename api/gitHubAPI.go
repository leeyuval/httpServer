package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const Organization = "1"
const Owner = "2"

const GitHubBaseUrl = "https://api.github.com/"
const organizationFilterExt = "orgs/%s/repos"
const ownertFilterExt = "users/%s/repos"

// GitHubAPI implements the API interface for GitHub repositories
type GitHubAPI struct{}

type githubRepoJson struct {
	Name  string `json:"name"`
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
	URL          string    `json:"html_url"`
	CreationTime time.Time `json:"created_at"`
	Stars        int       `json:"stargazers_count"`
}

func (api *GitHubAPI) BuildUrl(filter string, content string) string {
	var urlString = GitHubBaseUrl
	switch filter {
	case Organization:
		urlString += organizationFilterExt
	case Owner:
		urlString += ownertFilterExt
	default:
		urlString += organizationFilterExt
	}
	url := fmt.Sprintf(urlString, content)
	return url
}

func (api *GitHubAPI) DisplayResponse(response *http.Response) {
	var githubResponse []githubRepoJson

	err := json.NewDecoder(response.Body).Decode(&githubResponse)
	if err != nil {
		println(err)
	}

	fmt.Println("Fetched repositories:")
	for _, repo := range githubResponse {
		fmt.Printf("Name: %s\n", repo.Name)
		fmt.Printf("Owner: %s\n", repo.Owner.Login)
		fmt.Printf("URL: %s\n", repo.URL)
		fmt.Printf("Creation Time: %s\n", repo.CreationTime.String())
		fmt.Printf("Stars: %d\n", repo.Stars)
		fmt.Println("------")
	}
}

func (api *GitHubAPI) SendRequest(url string) *http.Response {
	response, err := http.Get(url)
	if err != nil {
		println(err)
	}
	return response
}

func (api *GitHubAPI) FetchRepositoriesByFilter(filter string, content string) {
	url := api.BuildUrl(filter, content)

	response := api.SendRequest(url)

	api.DisplayResponse(response)
}