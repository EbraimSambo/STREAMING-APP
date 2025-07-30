package core

type PageInfo struct {
	Total       int  `json:"total"`
	Page        int  `json:"page"`
	PerPage     int  `json:"perPage"`
	TotalPages  int  `json:"totalPages"`
	PrevPage    *int `json:"prevPage"`              // ponteiro para permitir null
	NextPage    *int `json:"nextPage"`              // ponteiro para permitir null
	HasNextPage bool `json:"hasNextPage,omitempty"` // opcional
}

// genérico para qualquer tipo T (Go 1.18+)
type Pagination[T any] struct {
	Items []*T     `json:"items"` // <-- aqui!
	Info  PageInfo `json:"info"`
}

type DataPagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}
