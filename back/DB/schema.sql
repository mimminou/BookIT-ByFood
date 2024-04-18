-- Using Sqlite,
-- Since SQLite does not (at least not directly) support length capping, I'm using a TEXT field, a VARCHAR of specified length would be used on another SQL DB

CREATE TABLE IF NOT EXISTS Books (
    book_id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    num_pages INTEGER,
    publication_date DATE NOT NULL,
);
