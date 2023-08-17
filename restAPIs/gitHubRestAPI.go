package restAPIs

import (
	"encoding/json"
	"fmt"
	"httpServer/requestsHandlers"
	"net/http"
	"time"
)

const Organization = "1"
const Owner = "2"

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

func (api *GitHubRestAPI) DisplayResponse(response *http.Response) {
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

func (api *GitHubRestAPI) FetchRepositoriesByFilter(filter string, content string) {
	url := api.BuildUrl(filter, content)

	response := api.SendGetRequest(url)

	api.DisplayResponse(response)
}
