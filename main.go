package main

import (
	"bufio"
	"fmt"
	"httpServer/restAPIs"
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

	var gitHubAPI restAPIs.GitHubRestAPI

	gitHubAPI.FetchRepositoriesByFilter(filter, content)

}
