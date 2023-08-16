package model

type API interface {
	FetchRepositoriesByFilter(filter string, content string)
}
