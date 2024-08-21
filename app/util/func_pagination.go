package util

import (
	"net/url"
	"strconv"
	"strings"
)

const (
	MAX_PAGE_SIZE = 99999999
)

type Pagination struct {
	Page      int64  `json:"page"`
	PageSize  int64  `json:"pageSize"`
	SortField string `json:"sortField"`
	SortOrder string `json:"sortOrder"`
}

type PaginationResponse[T any] struct {
	Count int64 `json:"count"`
	Rows  *[]T  `json:"rows"`
}

func (p *PaginationResponse[T]) Build() *PaginationResponse[T] {
	if len(*p.Rows) == 0 {
		emptyRows := make([]T, 0)
		p.Rows = &emptyRows
	}
	return p
}

func GetPageInQuery(query url.Values) (*Pagination, error) {
	pagination := Pagination{}
	page, err := strconv.ParseInt(query.Get("page"), 10, 64)
	if err != nil || query.Get("page") == "" {
		page = 1
	}
	pageSize, err := strconv.ParseInt(query.Get("pageSize"), 10, 64)
	if err != nil || query.Get("pageSize") == "" {
		pageSize = MAX_PAGE_SIZE
	}
	sortField := query.Get("sortField")
	if sortField == "" {
		sortField = "id"
	}
	sortOrder := query.Get("sortOrder")
	if sortOrder == "" {
		sortOrder = "desc"
	}
	pagination.Page = page
	pagination.PageSize = pageSize
	pagination.SortField = sortField
	pagination.SortOrder = sortOrder
	return &pagination, nil
}
func IsSortFieldOrder(column string, order string, sortField string, sortOrder string) bool {
	if strings.ToLower(column) == strings.ToLower(sortField) && strings.ToLower(order) == strings.ToLower(sortOrder) {
		return true
	}
	return false
}
