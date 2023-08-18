package main

import (
	"httpServer/restAPIs"
)

func main() {
	var gitHubAPI restAPIs.GitHubRestAPI
	gitHubAPI.FetchRepositoriesByType()
}
