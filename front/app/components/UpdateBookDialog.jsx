import { Dialog, DialogHeader, DialogContent, DialogTitle } from "@/components/ui/dialog"
import { useState } from "react"
import { Context } from '../context'
import { useContext } from 'react'


export default function UpdateBookDialog({ open, setOpen }) {
    const ctx = useContext(Context)
    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>
                        <UpdateBookForm books={ctx.books} setBooks={ctx.setBooks} setOpen={setOpen} book={ctx.selectedBook} book_id={ctx.selectedBookID} />
                    </DialogTitle>
                </DialogHeader>
            </DialogContent>
        </Dialog>
    )
}

function UpdateBookForm({ books, setBooks, setOpen, book, book_id }) {
    const selectedBook = book

    const [formData, setFormData] = useState({
        title: selectedBook.title,
        author: selectedBook.author,
        publicationDate: selectedBook.pub_date,
        totalPages: selectedBook.num_pages,
    });

    const [err, setErr] = useState({ title: false, author: false, publicationDate: false, totalPages: false })

    const handleChange = (event) => {
        const { name, value } = event.target;
        setFormData((prevData) => ({ ...prevData, [name]: value }));

        setErr((prevErrors) => {
            const newErrors = { ...prevErrors };
            switch (name) {
                case 'title':
                    newErrors[name] = value === "" ? true : false;
                    break;
                case 'author':
                    newErrors[name] = value === "" ? true : false;
                    break;
                case 'publicationDate':
                    newErrors[name] = !/^\d{4}-\d{2}-\d{2}$/.test(value) ? true : false;
                    break;
                case 'totalPages':
                    newErrors[name] = !/^\d+$/.test(value) ? true : false;
                    break;
                default:
                    break;
            }
            return newErrors;
        });
    };

    const handleSubmit = (event) => {
        event.preventDefault();

        // Validation before submission
        const errorList = []
        Object.entries(err).forEach(([key, isTrue]) => {
            if (isTrue) {
                errorList.push(key)
            }
        })
        if (errorList.length > 0) {
            alert('Please check the following fields : \n' + errorList);
            return;
        }
        //Optimistic update
        const newBook = {
            title: formData.title,
            author: formData.author,
            pub_date: formData.publicationDate,
            num_pages: parseInt(formData.totalPages),
        }
        selectedBook.book_id = book_id
        MakeRequest(newBook, book_id)


        // Submit form data (replace with your submission logic)
        console.log('Submitting form data:', formData);

        const currentBooks = books
        const updatedBookIndex = currentBooks.findIndex((book) => book.book_id === book_id)
        currentBooks[updatedBookIndex] = newBook
        setOpen(false)
        setBooks([...currentBooks])
    };

    return (
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            <label htmlFor="title">Book Title:</label>
            <input
                type="text"
                name="title"
                id="title"
                value={formData.title}
                onChange={handleChange}
                style={{ border: "1px solid gray" }}
                required
            />

            <label htmlFor="author">Author:</label>
            <input
                type="text"
                name="author"
                id="author"
                value={formData.author}
                onChange={handleChange}
                style={{ border: "1px solid gray" }}
                required
            />

            <label htmlFor="publicationDate">Year of Publication (YYYY-MM-DD):</label>
            <input
                type="text"
                name="publicationDate"
                id="publicationDate"
                value={formData.publicationDate}
                onChange={handleChange}
                pattern="^\d{4}-\d{2}-\d{2}$"
                style={{ border: "1px solid gray" }}
                required
            />

            <label htmlFor="totalPages">Number of Pages (Optional):</label>
            <input
                type="number"
                name="totalPages"
                id="totalPages"
                value={formData.totalPages}
                style={{ border: "1px solid gray" }}
                onChange={handleChange}
            />

            <button type="submit">Update</button>
        </form>
    );
}

function MakeRequest(Book, book_id) {

    fetch(`http://localhost:8046/books/${book_id}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(Book)
    })
}
