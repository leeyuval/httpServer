package requestsHandlers

import "net/http"

type RequestsHandler interface {
	SendGetRequest(url string) http.Response
}
