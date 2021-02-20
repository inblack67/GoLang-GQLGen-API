package mymodels

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// Story ...
type Story struct{
	gorm.Model
	ID uint `json:"-"`
	UUID  uuid.UUID `gorm:"unique; not null; default: null;" json:"uuid,omitempty"`
	Title string `json:"title,omitempty" gorm:"unique;not null;default:null"`
	UserID uint `json:"-" gorm:"not null"`
	UserUUID string `json:"userId" gorm:"not null"`
}

// ValidateStory ...
func (newStory Story) ValidateStory() error{
	return validation.ValidateStruct(&newStory, 

		validation.Field(&newStory.Title, validation.Required, validation.Length(5,20)),
	)
}