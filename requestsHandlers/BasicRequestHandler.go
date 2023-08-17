package requestsHandlers

import "net/http"

type BasicRequestsHandler struct{}

func (requestHandler *BasicRequestsHandler) SendGetRequest(url string) *http.Response {
	response, err := http.Get(url)
	if err != nil {
		println(err)
	}
	return response
}
