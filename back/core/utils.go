package core

import "os"

//Utility functions

// Writes anything passed to it to query.txt, This is only used for debugging
func LogQuery(query string) {
	file, _ := os.OpenFile("query.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	file.WriteString(query)
	defer file.Close()
}

// Verifiy if fields are empty, return empty fields
// Num pages is not checked because I decided to allow it to be empty

func CheckEmptyFields(book Book) []string {
	var emptyFields []string
	if book.Title == "" {
		emptyFields = append(emptyFields, "title")
	}
	if book.Author == "" {
		emptyFields = append(emptyFields, "author")
	}
	if book.Pub_Date == "" {
		emptyFields = append(emptyFields, "pub_date")
	}
	return emptyFields
}
