package usecase

type PaginationHeader struct {
	CurrentPage int64 `json:"current_page"`
	PerPage     int64 `json:"per_page"`
	TotalData   int64 `json:"total_data"`
	TotalPages  int64 `json:"total_pages"`
}
