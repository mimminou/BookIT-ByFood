package core

import (
	"database/sql"
)

/**
This file defines basic CRUD ops on DB
Refrained from using an ORM although I'm aware it's almost always better to use them, to showcase that I can use RAW SQL
All queries are done using Prepared Statements to directly mitigate SQL injections
**/

// get all books
func GetBooks(db *sql.DB) ([]Book, error) {
	rows, err := db.Query("SELECT * FROM Books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := make([]Book, 0)
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.Book_Id, &book.Title, &book.Author, &book.Num_Pages, &book.Pub_Date); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

// get single book
func GetBook(db *sql.DB, Book_id int) (Book, error) {
	var book Book
	err := db.QueryRow("SELECT * FROM Books WHERE Book_id = ?", Book_id).Scan(&book.Book_Id, &book.Title, &book.Author, &book.Num_Pages, &book.Pub_Date)
	return book, err
}

// Needed for pagination
func GetOffsetBooks(db *sql.DB, limit int, offset int) ([]Book, error) {
	rows, err := db.Query("SELECT * FROM Books LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := make([]Book, 0)
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.Book_Id, &book.Title, &book.Author, &book.Num_Pages, &book.Pub_Date); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

// add book
func AddBook(db *sql.DB, book Book) error {
	_, err := db.Exec("INSERT INTO Books (title, author, num_pages, pub_date) VALUES (?, ?, ?, ?)", book.Title, book.Author, book.Num_Pages, book.Pub_Date)
	return err
}

// delete book
func DeleteBook(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM Books WHERE book_id = ?", id)
	return err
}

// update book
// PUT request, not PATCH, so no need to do partial update
func UpdateBook(db *sql.DB, book Book) error {
	_, err := db.Exec("UPDATE Books SET title = ?, author = ?, num_pages = ?, pub_date = ? WHERE Book_id = ?", book.Title, book.Author, book.Num_Pages, book.Pub_Date, book.Book_Id)
	return err
}
