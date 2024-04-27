import { Dialog, DialogHeader, DialogContent, DialogTitle } from "@/components/ui/dialog"
import { Context } from '@/app/context'
import { Book, ErrMessage } from '@/app/types'
import { toast } from '@/components/ui/use-toast'
import { FormEvent, ChangeEvent, useContext, useState, Dispatch, SetStateAction } from 'react'

interface UpdateBookFormProps {
    books: Book[]
    setBooks: Dispatch<SetStateAction<Book[]>>
    setOpen: Dispatch<SetStateAction<boolean>>
    book: Book
    toaster: typeof toast
    setBook?: Dispatch<SetStateAction<Book>> //Optional field,
}

interface UpdateBookDialogProps {
    open: boolean
    setOpen: Dispatch<SetStateAction<boolean>>
    setBook?: Dispatch<SetStateAction<Book>>
}


interface RequestProps {
    oldBook: Book
    newBook: Book
    books: Book[]
    setBooks: Dispatch<SetStateAction<Book[]>>
    toaster: typeof toast
    setBook?: Dispatch<SetStateAction<Book>> //Optional field, needed to update state of the book in the book detail page
}


export default function UpdateBookDialog({ open, setOpen, setBook }: UpdateBookDialogProps) {
    const ctx = useContext(Context)
    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>
                        Update Book
                    </DialogTitle>
                </DialogHeader>
                {setBook ?
                    <UpdateBookForm books={ctx.books} setBooks={ctx.setBooks} setBook={setBook} setOpen={setOpen} book={ctx.selectedBook} toaster={ctx.toast} />
                    :
                    <UpdateBookForm books={ctx.books} setBooks={ctx.setBooks} setOpen={setOpen} book={ctx.selectedBook} toaster={ctx.toast} />
                }
            </DialogContent>
        </Dialog>
    )
}

function UpdateBookForm(props: UpdateBookFormProps) {
    const { books, setBooks, setBook, setOpen, book, toaster } = props

    const oldBook = book
    const [formData, setFormData] = useState({
        title: book.title,
        author: book.author,
        publicationDate: book.pub_date,
        totalPages: book.num_pages === 0 ? "" : book.num_pages,
    });

    const [err, setErr] = useState({ title: false, author: false, publicationDate: false, totalPages: false })

    const isDateValid = (date: string) => {
        const d = new Date(date)
        const dNum = d.getTime()
        if (!dNum && dNum !== 0) {
            return false
        }
        return true
    }

    const handleChange = (event: ChangeEvent<HTMLInputElement>) => {
        const { name, value } = event.target;
        if (name == "totalPages") {
            console.log(name, value, typeof value)
        }
        console.log(name, value)
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

    const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
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

        console.log(formData)

        //Optimistic update
        const newBook = {
            book_id: book.book_id,
            title: formData.title,
            author: formData.author,
            pub_date: formData.publicationDate,
            num_pages: parseInt(formData.totalPages.toString()), // why do I need toString() here ? Wihtout it, TS complains
        }

        MakeRequest({ newBook, oldBook, books, setBooks, setBook, toaster })

        const updatedBook = {
            book_id: book.book_id,
            title: formData.title,
            author: formData.author,
            pub_date: formData.publicationDate,
            num_pages: formData.totalPages,
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

async function MakeRequest(props: RequestProps) {
    const { oldBook, newBook, books, setBooks, toaster, setBook } = props
    try {
        const response = await fetch(`http://localhost:8046/books/${oldBook.book_id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(newBook)
        })
        let jsonresponse = await response.json()
        if (!response.ok) {
            // Throw on Failure, whether we get a failure message or not
            jsonresponse = jsonresponse as ErrMessage
            if ("msg" in jsonresponse) {
                throw new Error(jsonresponse)
            }
            else {
                throw new Error("Unkown response type")
            }
        }
        toaster({
            description: 'Book Updated',
        })
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
        if ("msg" in error) {
            toaster({
                title: 'Operation Failed',
                description: error.msg,
                variant: 'destructive',
            })
            return
        }
        toaster({
            title: 'Operation Failed',
            description: 'Network Error, Could not connect to server',
            variant: 'destructive',
        })
    }
}
