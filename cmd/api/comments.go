// Filename: cmd/api/comments.go
package main

import (
	"errors"
	"fmt"
	"net/http"

	// import the data package which contains the definition for Comment

	"github.com/Lee26Ed/qod/internal/data"
	"github.com/Lee26Ed/qod/internal/validator"
)

func (a *application) createCommentHandler(w http.ResponseWriter,
                                                      r *http.Request) { 
    // create a struct to hold a comment
    // we use struct tags[``] to make the names display in lowercase
    var incomingData struct {
        Content  string  `json:"content"`
        Author   string  `json:"author"`
    }

	// perform the decoding
   err := a.readJSON(w, r, &incomingData)
   if err != nil {
		a.badRequestResponse(w, r, err)
       return
   }

   // Copy the values from incomingData to a new Comment struct
	// At this point in our code the JSON is well-formed JSON so now
	// we will validate it using the Validator which expects a Comment
	quote := &data.Quotes {
		Content: incomingData.Content,
		Author: incomingData.Author,
	}
	// Initialize a Validator instance
	v := validator.New()

	data.ValidateQuote(v, quote)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Add the quote to the database table
   err = a.quoteModel.Insert(quote)
   if err != nil {
       a.serverErrorResponse(w, r, err)
       return
   }


// for now display the result
   fmt.Fprintf(w, "%+v\n", incomingData)
      // Set a Location header. The path to the newly created quote
   headers := make(http.Header)
   headers.Set("Location", fmt.Sprintf("/v1/quotes/%d", quote.ID))
   // Send a JSON response with 201 (new resource created) status code
  data := envelope{
         "quote": quote,
       }
  err = a.writeJSON(w, http.StatusCreated, data, headers)
  if err != nil {
       a.serverErrorResponse(w, r, err)
       return
   }


}

func (a *application)displayQuoteHandler(
                                               w http.ResponseWriter,
                                               r *http.Request) {

	// Get the id from the URL /v1/quotes/:id so that we
	// can use it to query the quotes table. We will 
	// implement the readIDParam() function later
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
       return 
	}

	// Call Get() to retrieve the quote with the specified id
	quote, err := a.quoteModel.Get(id)
	if err != nil {
		switch {
			case errors.Is(err, data.ErrRecordNotFound):
				a.notFoundResponse(w, r)
			default:
				a.serverErrorResponse(w, r, err)
		}
		return 
	}

	// display the quote
    data := envelope {
                "quote": quote,
            }
    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
       a.serverErrorResponse(w, r, err)
       return 
   }

}

