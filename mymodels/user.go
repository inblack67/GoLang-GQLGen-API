package mymodels

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"gorm.io/gorm"
)

// User ...
type User struct{
	gorm.Model
	ID uint `json:"-"`
	UUID  string `gorm:"unique; not null; default: null;" json:"uuid,omitempty"`
	Name string `json:"name,omitempty" gorm:"not null;default:null"`
	Username string `json:"username,omitempty" gorm:"unique;not null;default:null"`
	Email string `json:"email,omitempty" gorm:"unique;not null;default:null"`
	Password string `json:"-" gorm:"not null;default:null"`
	Stories []Story `json:"stories" gorm:"foreignKey:UserID"`
}

// ValidateMe ...
func (newUser User) ValidateMe() error{
	return validation.ValidateStruct(&newUser, 

		validation.Field(&newUser.Name, validation.Required, validation.Length(5,20)),

		validation.Field(&newUser.Username, validation.Required, validation.Length(5,20)),

		validation.Field(&newUser.Email, validation.Required, is.Email),

		validation.Field(&newUser.Password, validation.Required, validation.Length(8,20)),
	)
}
