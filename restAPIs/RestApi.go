package restAPIs

import "net/http"

type RestAPI interface {
	GetUserInput() (filter string, content string, err error)
	BuildUrl(filter string, content string) string
	SendGetRequest(url string) http.Response
	DisplayResponse(response string)
	FetchRepositoriesByFilter(filter string, content string)
}
