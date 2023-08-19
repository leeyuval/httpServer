package restAPIs

type RestAPI interface {
	GetUserInput() (filter string, content string, err error)
	BuildUrl(filter string, content string) string
	DisplayResponse(response string)
	FetchRepositoriesByType(filter string, content string)
}
