package restAPIs

import "net/http"

type RestAPI interface {
	BuildUrl(filter string, content string) string
	SendGetRequest(url string) http.Response
	DisplayResponse(response string)
	FetchRepositoriesByFilter(filter string, content string)
}
