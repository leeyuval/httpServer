package restAPIs

import (
	"encoding/json"
	"fmt"
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

func (api *GitHubRestAPI) DisplayResponse(response *http.Response) ([]githubRepoJson, error) {
	var githubResponse []githubRepoJson

	err := json.NewDecoder(response.Body).Decode(&githubResponse)
	if err != nil {
		return nil, err
	}

	return githubResponse, nil
}

func (api *GitHubRestAPI) FetchRepositoriesByFilter(filter string, content string) ([]githubRepoJson, error) {
	url := api.BuildUrl(filter, content)

	response := api.SendGetRequest(url)

	repositories, err := api.DisplayResponse(response)
	if err != nil {
		return nil, err
	}

	return repositories, nil
}
