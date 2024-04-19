package server

import (
	"database/sql"
	models "github.com/mimminou/BookIT-ByFood/back/models"
)

/**
This file defines basic CRUD ops on DB
Refrained from using an ORM although I'm aware it's almost always better to use them, to showcase that I can use RAW SQL
All queries are done using Prepared Statements to directly mitigate SQL injections
**/

// get all books
func GetBooks(db *sql.DB) ([]models.Book, error) {
	rows, err := db.Query("SELECT * FROM Books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := make([]models.Book, 0)
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.Book_Id, &book.Title, &book.Author, &book.Num_Pages, &book.Pub_Date); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

// get single book
func GetBook(db *sql.DB, Book_id int) (models.Book, error) {
	var book models.Book
	err := db.QueryRow("SELECT * FROM Books WHERE Book_id = ?", Book_id).Scan(&book.Book_Id, &book.Title, &book.Author, &book.Num_Pages, &book.Pub_Date)
	return book, err
}

// add book
func AddBook(db *sql.DB, book models.Book) (int, error) {
	operation, err := db.Exec("INSERT INTO Books (title, author, num_pages, pub_date) VALUES (?, ?, ?, ?)", book.Title, book.Author, book.Num_Pages, book.Pub_Date)
	if err != nil {
		return 0, err
	}
	RowsInserted, err := operation.LastInsertId()
	if RowsInserted == 0 {
		return 0, sql.ErrNoRows
	}

	return int(RowsInserted), err
}

// delete book
func DeleteBook(db *sql.DB, id int) error {
	operation, err := db.Exec("DELETE FROM Books WHERE book_id = ?", id)
	if err != nil {
		return err
	}
	RowsDeleted, err := operation.RowsAffected()
	if RowsDeleted == 0 {
		return sql.ErrNoRows
	}
	return err
}

// update book
// PUT request, not PATCH, so no need to do partial update
func UpdateBook(db *sql.DB, book models.Book) error {
	operation, err := db.Exec("UPDATE Books SET title = ?, author = ?, num_pages = ?, pub_date = ? WHERE Book_id = ?", book.Title, book.Author, book.Num_Pages, book.Pub_Date, book.Book_Id)
	if err != nil {
		return err
	}
	RowsUpdated, err := operation.RowsAffected()
	if RowsUpdated == 0 {
		return sql.ErrNoRows
	}
	return err
}
