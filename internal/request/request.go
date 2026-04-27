package request

type PaginationRequest struct {
	AfterCursor string `query:"after" validate:"omitempty,uuid4"`
	Limit       int    `query:"limit" validate:"omitempty,number" default:"50"`
}
