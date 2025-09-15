// Filename: internal/data/comments.go
package data

import (
	"time"

	"github.com/Lee26Ed/qod/internal/validator"
)

// each name begins with uppercase so that they are exportable/public
type Quotes struct {
    ID int64                     `json:"id"`                   
    Content  string              `json:"content"`     
    Author  string               `json:"author"`
    CreatedAt  time.Time         `json:"-"`     
    Version int32                `json:"version"`      
} 


// Create a function that performs the validation checks
func ValidateComment(v *validator.Validator, quote *Quotes) {
	// check if the Content field is empty
    v.Check(quote.Content != "", "content", "must be provided")
	// check if the Author field is empty
    v.Check(quote.Author != "", "author", "must be provided")
	// check if the Content field is empty
    v.Check(len(quote.Content) <= 100, "content", "must not be more than 100 bytes long")
	// check if the Author field is empty
     v.Check(len(quote.Author) <= 25, "author", "must not be more than 25 bytes long")
}