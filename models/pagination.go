package models

type Pagination struct {
	Next string `json:"next"`
	Prev string `json:"prev"`
}

type Response[T any] struct {
	Data T `json:"data"`
}

type WithPagination[T any] struct {
	Data   T          `json:"data"`
	Paging Pagination `json:"paging"`
}
