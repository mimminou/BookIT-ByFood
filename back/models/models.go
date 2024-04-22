package models

//Here are the structs that define the models of the DB, needed for unmarshalling/Marhsalling Json

// Book is the schema for a book object
// @Description Book
type Book struct {
	// @Property book_id int true "Book ID"
	Book_Id int `json:"book_id"`
	// @Property title string true "Title"
	Title string `json:"title"`
	// @Property author string true "Author"
	Author string `json:"author"`
	// @Property num_pages string false "Number of pages"
	Num_Pages *int `json:"num_pages,omitempty"`
	// @Property pub_date int true "Publication date"
	Pub_Date string `json:"pub_date"`
}
