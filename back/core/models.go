package core

//Here are the structs that define the models of the DB, needed for unmarshalling/Marhsalling Json

type Book struct {
	Book_Id   int    `json:"book_id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Num_Pages *int   `json:"num_pages,omitempty"`
	Pub_Date  string `json:"pub_date"`
}
