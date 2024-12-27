package dto

import (
	"encoding/json"
	"errors"
	"github.com/lyh-demo/go-webapp-demo/model"
	"gopkg.in/go-playground/validator.v9"
)

const (
	required string = "required"
	maxVal   string = "max"
	minVal   string = "min"
)

// BookDto defines a data transfer object for book.
type BookDto struct {
	Title      string `validate:"required,min=3,max=50" json:"title"`
	Isbn       string `validate:"required,min=10,max=20" json:"isbn"`
	CategoryID uint   `json:"categoryId"`
	FormatID   uint   `json:"formatId"`
	messages   map[string]string
}

// NewBookDto is constructor.
func NewBookDto(messages map[string]string) *BookDto {
	return &BookDto{messages: messages}
}

// Create creates a book model from this DTO.
func (b *BookDto) Create() *model.Book {
	return model.NewBook(b.Title, b.Isbn, b.CategoryID, b.FormatID)
}

// Validate performs validation check for the item.
func (b *BookDto) Validate() map[string]string {
	return validateDto(b)
}

func validateDto(b *BookDto) map[string]string {
	err := validator.New().Struct(b)
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	errors.As(err, &validationErrors)
	if len(validationErrors) == 0 {
		return nil
	}

	return createErrorMessages(b, validationErrors)
}

func createErrorMessages(b *BookDto, errors validator.ValidationErrors) map[string]string {
	result := make(map[string]string)
	for i := range errors {
		switch errors[i].StructField() {
		case "Title":
			switch errors[i].Tag() {
			case required, minVal, maxVal:
				result["title"] = b.messages["ValidationErrMessageBookTitle"]
			}
		case "Isbn":
			switch errors[i].Tag() {
			case required, minVal, maxVal:
				result["isbn"] = b.messages["ValidationErrMessageBookISBN"]
			}
		}
	}
	return result
}

// ToString is return string of object
func (b *BookDto) ToString() (string, error) {
	bytes, err := json.Marshal(b)
	return string(bytes), err
}
