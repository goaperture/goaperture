package client

import (
	"context"

	"github.com/goaperture/goaperture/v2/responce"
)

type PaginationTotalKey struct{}

type RequestPagination struct {
	Use bool
	responce.Pagination
}

func WithPagination(ctx context.Context) context.Context {
	return context.WithValue(ctx, PaginationTotalKey{}, new(RequestPagination))
}

func SetPagination(ctx context.Context, pagination responce.Pagination) {
	request := GetPagination(ctx)

	request.Pagination = pagination
	request.Use = true
}

func GetPagination(ctx context.Context) *RequestPagination {
	val := ctx.Value(PaginationTotalKey{})
	if req, ok := val.(*RequestPagination); ok {
		return req
	}
	return nil
}

func (p *RequestPagination) Export() *responce.Pagination {
	if !p.Use {
		return nil
	}

	return &p.Pagination
}
