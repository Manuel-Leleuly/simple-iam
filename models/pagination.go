package models

type Pagination struct {
	Next string `json:"next"`
	Prev string `json:"prev"`
}

type WithPagination struct {
	Data   any        `json:"data"`
	Paging Pagination `json:"paging"`
}
