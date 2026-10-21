package aperture

import "github.com/goaperture/goaperture/v3/response"

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type Response struct {
	Data       any                  `json:"data"`
	Error      *Error               `json:"error,omitempty"`
	Pagination *response.Pagination `json:"pagination,omitempty"`
}
