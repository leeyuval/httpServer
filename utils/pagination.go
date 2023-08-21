package utils

// PaginationData holds pagination-related information.
type PaginationData struct {
	Page       int
	PrevPage   int
	NextPage   int
	HasPrev    bool
	HasNext    bool
	TotalPages int
}

// GetPaginationData calculates pagination-related data based on the current page and total count.
func GetPaginationData(currentPage, itemsPerPage, totalCount int) PaginationData {
	totalPages := (totalCount + itemsPerPage - 1) / itemsPerPage

	return PaginationData{
		Page:       currentPage,
		PrevPage:   currentPage - 1,
		NextPage:   currentPage + 1,
		HasPrev:    currentPage > 1,
		HasNext:    currentPage < totalPages,
		TotalPages: totalPages,
	}
}
