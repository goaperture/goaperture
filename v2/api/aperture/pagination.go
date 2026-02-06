package aperture

import "context"

type Pagination struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

func NewPagination(ctx context.Context, page int, size int) Pagination {
	return Pagination{Page: page, Size: size, Total: -1}
}
