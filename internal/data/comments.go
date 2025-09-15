// Filename: internal/data/comments.go
package data

import (
	"context"
	"database/sql"
	"errors"
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

// A QuoteModel expects a connection pool
type QuoteModel struct {
    DB *sql.DB
}

// Create a function that performs the validation checks
func ValidateQuote(v *validator.Validator, quote *Quotes) {
	// check if the Content field is empty
    v.Check(quote.Content != "", "content", "must be provided")
	// check if the Author field is empty
    v.Check(quote.Author != "", "author", "must be provided")
	// check if the Content field is empty
    v.Check(len(quote.Content) <= 100, "content", "must not be more than 100 bytes long")
	// check if the Author field is empty
     v.Check(len(quote.Author) <= 25, "author", "must not be more than 25 bytes long")
}

// Insert a new row in the quotes table
// Expects a pointer to the actual quote
func (q QuoteModel) Insert(quote *Quotes) error {
   // the SQL query to be executed against the database table
    query := `
        INSERT INTO quotes (content, author)
        VALUES ($1, $2)
        RETURNING id, created_at, version
        `
  // the actual values to replace $1, and $2
   args := []any{quote.Content, quote.Author}
 
	// Create a context with a 3-second timeout. No database
	// operation should take more than 3 seconds or we will quit it
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()
	// execute the query against the quotes database table. We ask for the the
	// id, created_at, and version to be sent back to us which we will use
	// to update the Quote struct later on
	return q.DB.QueryRowContext(ctx, query, args...).Scan(
														&quote.ID,
														&quote.CreatedAt,
														&quote.Version)

}


// Get a specific Quote from the quotes table
func (q QuoteModel) Get(id int64) (*Quotes, error) {
   // check if the id is valid
    if id < 1 {
        return nil, ErrRecordNotFound
    }
   // the SQL query to be executed against the database table
    query := `
        SELECT id, created_at, content, author, version
        FROM quotes
        WHERE id = $1
      `
	// declare a variable of type Quote to store the returned quote
	var quote Quotes

	// Set a 3-second context/timer
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	err := q.DB.QueryRowContext(ctx, query, id).Scan (
												&quote.ID,
												&quote.CreatedAt,
												&quote.Content,
												&quote.Author,
												&quote.Version,
												)
	// check for which type of error
	if err != nil {
		switch {
			case errors.Is(err, sql.ErrNoRows):
				return nil, ErrRecordNotFound
			default:
				return nil, err
			}
		}
	return &quote, nil
	}

// Update a specific Quote from the quotes table
func (q QuoteModel) Update(quote *Quotes) error {
	// The SQL query to be executed against the database table
	// Every time we make an update, we increment the version number
	query := `
        UPDATE quotes
        SET content = $1, author = $2, version = version + 1
        WHERE id = $3
        RETURNING version
      `
   args := []any{quote.Content, quote.Author, quote.ID}
   ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
   defer cancel()

   return q.DB.QueryRowContext(ctx, query, args...).Scan(&quote.Version)

}

// Delete a specific Quote from the quotes table
func (q QuoteModel) Delete(id int64) error {

    // check if the id is valid
    if id < 1 {
        return ErrRecordNotFound
    }
   // the SQL query to be executed against the database table
    query := `
        DELETE FROM quotes
        WHERE id = $1
      `

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
   defer cancel()

// ExecContext does not return any rows unlike QueryRowContext. 
// It only returns  information about the the query execution
// such as how many rows were affected
   result, err := q.DB.ExecContext(ctx, query, id)
   if err != nil {
       return err
   }

   // Were any rows  delete?
    rowsAffected, err := result.RowsAffected()
    if err != nil {
       return err
   }
// Probably a wrong id was provided or the client is trying to
// delete an already deleted quote
   if rowsAffected == 0 {
       return ErrRecordNotFound
   }

   return nil

}

// Get all comments
func (q QuoteModel) GetAll() ([]*Quotes, error) {

	// the SQL query to be executed against the database table
    query := `
        SELECT id, created_at, content, author, version
        FROM quotes
        ORDER BY id
      `
   ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
   defer cancel()

   // QueryContext returns multiple rows.
	rows, err := q.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	// clean up the memory that was used
	defer rows.Close()
	// we will store the address of each comment in our slice
	quotes := []*Quotes{}

	// process each row that is in rows

	for rows.Next() {
		var quote Quotes
		err := rows.Scan(&quote.ID,
						&quote.CreatedAt,
						&quote.Content,
						&quote.Author,
						&quote.Version,
						)
		if err != nil {
			return nil, err
		}
	// add the row to our slice
	quotes = append(quotes, &quote)
	}  // end of for loop

	// after we exit the loop we need to check if it generated any errors
	err = rows.Err()
	if err != nil {
	return nil, err
	}

	return quotes, nil

}
