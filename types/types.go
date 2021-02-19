package types

import "net/http"

// SSession ...
type SSession struct{
	Username string
	ID uint 
	UUID string
}

// Status ...
type Status struct{
	Success bool `json:"success"`
	Message string `json:"message"`
}

// MyCtx ...
type MyCtx struct {
	Request *http.Request
	ResponseWriter http.ResponseWriter
}