package model

type Response struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    interface{}   `json:"data,omitempty"`
	Meta    *PaginatedMeta `json:"meta,omitempty"`
}
