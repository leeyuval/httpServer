package model

import "net/http"

type API interface {
	BuildUrl(filter string, content string) string
	SendRequest(url string) http.Response
	DisplayResponse(response string)
	FetchRepositoriesByFilter(filter string, content string)
}
