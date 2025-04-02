package model

type Response struct {
	Message string         `json:"message"`
	Data    any            `json:"data,omitempty"`
	Meta    *PaginatedMeta `json:"meta,omitempty"`
}

type FailedResponse struct {
	Message string `json:"message"`
	Errors  any    `json:"errors"`
}
