package aperture

import "github.com/goaperture/goaperture/v2/responce"

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type Responce struct {
	Data       any                  `json:"data"`
	Error      *Error               `json:"error,omitempty"`
	Pagination *responce.Pagination `json:"pagination,omitempty"`
}
