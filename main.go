package main

import (
	"bufio"
	"fmt"
	"httpServer/api"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please provide a filter (1 - Organization, 2 - Owner):")
	filter, _ := reader.ReadString('\n')
	filter = strings.TrimSpace(filter)
	fmt.Println("Please provide an appropriate content:")
	content, _ := reader.ReadString('\n')
	content = strings.TrimSpace(content)

	var gitHubAPI api.GitHubAPI

	url := gitHubAPI.BuildUrl(filter, content)

	response := gitHubAPI.SendRequest(url)

	gitHubAPI.DisplayResponse(response)

}
