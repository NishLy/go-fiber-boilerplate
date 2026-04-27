package response

type GenericResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type GenericSuccessResponse[T any] struct {
	GenericResponse
	Data T `json:"data"`
}

type PaginationMeta struct {
	Before  string `json:"before,omitempty"`
	After   string `json:"after,omitempty"`
	HasNext bool   `json:"has_next"`
	HasPrev bool   `json:"has_prev"`
}

type PagedDataResponse[T any] struct {
	GenericResponse
	Data []T            `json:"data"`
	Meta PaginationMeta `json:"meta"`
}
