package httpServer

// Repository struct to hold repository information
type Repository struct{}

// API interface defines the methods required by different APIs
type API interface {
	FetchRepositories(orgName string) ([]Repository, error)
}

// GitHubAPI implements the API interface for GitHub repositories
type GitHubAPI struct{}

func (api *GitHubAPI) FetchRepositories(orgName string) ([]Repository, error) {}

func main() {}
