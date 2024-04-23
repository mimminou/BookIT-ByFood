import { Dialog, DialogHeader, DialogContent, DialogTitle } from "@/components/ui/dialog"
import { useState } from "react"
import { Context } from '@/app/context'
import { useContext } from 'react'


export default function UpdateBookDialog({ open, setOpen, setBook }) {
    const ctx = useContext(Context)
    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>
                        Update Book
                    </DialogTitle>
                </DialogHeader>
                <UpdateBookForm books={ctx.books} setBooks={ctx.setBooks} setBook={setBook} setOpen={setOpen} book={ctx.selectedBook} toaster={ctx.toast} />
            </DialogContent>
        </Dialog>
    )
}

function UpdateBookForm({ books, setBooks, setBook, setOpen, book, toaster }) {

    const oldBook = book
    const [formData, setFormData] = useState({
        title: book.title,
        author: book.author,
        publicationDate: book.pub_date,
        totalPages: book.num_pages ? book.num_pages : '',
    });

    const [err, setErr] = useState({ title: false, author: false, publicationDate: false, totalPages: false })


    const isDateValid = (date) => {
        const d = new Date(date)
        const dNum = d.getTime()
        if (!dNum && dNum !== 0) {
            return false
        }
        return true
    }

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
                    newErrors[name] =
                        !/^\d{4}-\d{2}-\d{2}$/.test(value) ? true : !isDateValid(value)
                    break;
                case 'totalPages':
                    newErrors[name] = !/^\d+$/.test(value) && value !== "" ? true : false;
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
        const fieldNameMap = {
            title: 'Book Title',
            author: 'Author',
            publicationDate: 'Publication Date',
            totalPages: 'Number of Pages',
        }
        Object.entries(err).forEach(([key, isTrue]) => {
            if (isTrue) {
                errorList.push(key)
            }
        })
        if (errorList.length > 0) {
            const formattedErrList = errorList.map((key) => fieldNameMap[key])
            alert('Please check the following fields : \n' + formattedErrList.join("\n"));
            return;
        }

        //Optimistic update
        const newBook = {
            title: formData.title,
            author: formData.author,
            pub_date: formData.publicationDate,
            num_pages: parseInt(formData.totalPages),
        }

        MakeRequest(newBook, oldBook, books, setBooks, setBook, toaster)

        const updatedBook = {
            book_id: book.book_id,
            title: formData.title,
            author: formData.author,
            pub_date: formData.publicationDate,
            num_pages: parseInt(formData.totalPages)
        }
        if (setBook) {
            setBook(updatedBook)
        }

        const updatedBookList = books.map((element) => {
            if (element.book_id === updatedBook.book_id) {
                return updatedBook
            }
            return element
        })

        setOpen(false)
        setBooks(updatedBookList)

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

            <label htmlFor="publicationDate">Publication Date (YYYY-MM-DD):</label>
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

            <button type="submit" className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Update</button>
        </form>
    );
}


async function MakeRequest(newBook, oldBook, books, setBooks, setBook, toaster) {

    try {
        const response = await fetch(`http://localhost:8046/books/${oldBook.book_id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(newBook)
        })

        if (response.ok) {
            toaster({
                description: 'Book Updated',
            })
            return
        }
        else {
            const json = await response.json()
            if (setBook) {
                setBook(oldBook)
            }
            setBooks(books.map((element) => {
                if (element.book_id === oldBook.book_id) {
                    return oldBook
                }
                return element
            }))
            console.log(json.msg)
            toaster({
                title: 'Operation Failed',
                description: json.msg ? json.msg : 'Book Update Failed',
                variant: 'destructive',
            })
        }
    }
    catch (error) {
        if (setBook) {
            setBook(oldBook)
        }
        setBooks(books.map((element) => {
            if (element.book_id === oldBook.book_id) {
                return oldBook
            }
            return element
        }))
        toaster({
            title: 'Operation Failed',
            description: 'Network Error, Could not connect to server',
            variant: 'destructive',
        })
        return
    }
}
