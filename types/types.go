package types

// SSession ...
type SSession struct{
	Username string
	ID uint 
}

// Status ...
type Status struct{
	Success bool `json:"success"`
	Message string `json:"message"`
}