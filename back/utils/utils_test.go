package utils

import (
	models "github.com/mimminou/BookIT-ByFood/back/models"
	"testing"
)

func TestCheckEmptyFields(t *testing.T) {

	t.Run("Check with full Book struct", func(t *testing.T) {
		book := models.Book{
			Title:    "test",
			Author:   "test",
			Pub_Date: "test",
		}
		num := 42
		book.Num_Pages = &num

		emptyFields := CheckEmptyFields(book)
		if len(emptyFields) != 0 {
			t.Errorf("Expected empty fields array to be empty, got %v", emptyFields)
		}
	})

	t.Run("Check with empty Book struct", func(t *testing.T) {
		book := models.Book{}
		emptyFields := CheckEmptyFields(book)
		if len(emptyFields) == 0 {
			t.Errorf("Expected empty fields array to be not empty, got %v", emptyFields)
		}
	})

	t.Run("Check with empty Num_Pages field", func(t *testing.T) {
		book := models.Book{
			Title:    "test",
			Author:   "test",
			Pub_Date: "test",
		}

		emptyFields := CheckEmptyFields(book)
		if len(emptyFields) != 0 {
			t.Errorf("Expected empty fields array to be empty, got %v", emptyFields)
		}
	})

	t.Run("Check with a mandatory value missing (author)", func(t *testing.T) {
		book := models.Book{
			Title:    "test",
			Pub_Date: "test",
		}
		emptyFields := CheckEmptyFields(book)
		if len(emptyFields) == 0 {
			t.Errorf("Expected empty fields array to not be empty, got %v", emptyFields)
		}
	})

}
