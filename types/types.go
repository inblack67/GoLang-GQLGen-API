package types

import "github.com/inblack67/GQLGenAPI/mymodels"

// SSession ...
type SSession struct{
	User mymodels.User
}

// Status ...
type Status struct{
	Success bool `json:"success"`
	Message string `json:"message"`
}