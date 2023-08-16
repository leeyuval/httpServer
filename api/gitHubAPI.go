package api

import (
	"encoding/json"
	"fmt"
	"httpServer/model"
	"net/http"
	"time"
)

// GitHubAPI implements the API interface for GitHub repositories
type GitHubAPI struct{}

func (api *GitHubAPI) FetchRepositories(orgName string) ([]model.Repository, error) {
	url := fmt.Sprintf("https://api.github.com/orgs/%s/repos", orgName)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var githubResponse []struct {
		Name  string `json:"name"`
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
		HTMLURL         string    `json:"html_url"`
		CreatedAt       time.Time `json:"created_at"`
		StargazersCount int       `json:"stargazers_count"`
	}

	err = json.NewDecoder(response.Body).Decode(&githubResponse)
	if err != nil {
		return nil, err
	}

	var repositories []model.Repository
	for _, repo := range githubResponse {
		repositories = append(repositories, model.Repository{
			Name:         repo.Name,
			Owner:        repo.Owner,
			URL:          repo.HTMLURL,
			CreationTime: repo.CreatedAt,
			Stars:        repo.StargazersCount,
		})
		println(repo.Name)
	}

	return repositories, nil
}
