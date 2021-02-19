package mymodels

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// User ...
type User struct{
	gorm.Model
	ID uint `json:"-"`
	UUID  uuid.UUID `gorm:"unique; not null; default: null;" json:"uuid,omitempty"`
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

// Story ...
type Story struct{
	gorm.Model
	ID uint `json:"-"`
	UUID  uuid.UUID `gorm:"unique; not null; default: null;" json:"uuid,omitempty"`
	Title string `json:"title,omitempty" gorm:"unique;not null;default:null"`
	UserID uint `json:"-" gorm:"not null"`
}

// ValidateStory ...
func (newStory Story) ValidateStory() error{
	return validation.ValidateStruct(&newStory, 

		validation.Field(&newStory.Title, validation.Required, validation.Length(5,20)),
	)
}