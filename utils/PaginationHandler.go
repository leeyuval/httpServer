package paginationHandler

import (
	"math"
	"net/http"
	"strconv"
)

func GetCurrentPage(r *http.Request) int {
	pageParam := r.URL.Query().Get("page")
	currentPage, _ := strconv.Atoi(pageParam)
	if currentPage <= 0 {
		currentPage = 1
	}
	return currentPage
}

func GetPageIndices(currentPage int, perPage int, totalItems int) (int, int) {
	startIndex := (currentPage - 1) * perPage
	endIndex := startIndex + perPage
	if endIndex > totalItems {
		endIndex = totalItems
	}
	return startIndex, endIndex
}

func GetPageNumbers(currentPage int, totalPages int) []int {
	startPage := int(math.Max(float64(currentPage-2), 1))
	endPage := int(math.Min(float64(currentPage+2), float64(totalPages)))

	pageNumbers := make([]int, endPage-startPage+1)
	for i := startPage; i <= endPage; i++ {
		pageNumbers[i-startPage] = i
	}
	return pageNumbers
}

func GetQueryParams(r *http.Request) map[string]string {
	queryParams := make(map[string]string)
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			queryParams[key] = values[0]
		}
	}
	return queryParams
}
