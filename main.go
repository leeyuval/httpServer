package main

import (
	"bufio"
	"fmt"
	"httpServer/api"
	"log"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please provide an organization name:")
	orgName, _ := reader.ReadString('\n')
	orgName = strings.TrimSpace(orgName)

	if orgName == "" {
		fmt.Println("Organization name is required")
		return
	}

	var gitHubAPI api.GitHubAPI

	repositories, err := gitHubAPI.FetchRepositories(orgName)

	if err != nil {
		log.Fatal("Error fetching repositories:", err)
		return
	}

	fmt.Println("Fetched repositories:")
	for _, repo := range repositories {
		fmt.Printf("Name: %s\n", repo.Name)
		fmt.Printf("Owner: %s\n", repo.Owner)
		fmt.Printf("URL: %s\n", repo.URL)
		fmt.Printf("Creation Time: %s\n", repo.CreationTime.String())
		fmt.Printf("Stars: %d\n", repo.Stars)
		fmt.Println("------")
	}

}
