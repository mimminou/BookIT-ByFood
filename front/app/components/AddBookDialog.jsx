import { Dialog, DialogHeader, DialogContent, DialogTitle } from "@/components/ui/dialog"
import { useContext } from "react"
import { Context } from '../context'
import { useState } from "react"


export default function AddBookDialog({ open, setOpen }) {
    const ctx = useContext(Context)
    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>
                        <AddBookForm books={ctx.books} setBooks={ctx.setBooks} setOpen={setOpen} toaster={ctx.toast} />
                    </DialogTitle>
                </DialogHeader>
            </DialogContent>
        </Dialog>
    )
}

function AddBookForm({ books, setBooks, setOpen, toaster }) {

    const [formData, setFormData] = useState({
        title: '',
        author: '',
        publicationDate: '',
        totalPages: '',
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

        const newBook = {
            title: formData.title,
            author: formData.author,
            pub_date: formData.publicationDate,
            num_pages: parseInt(formData.totalPages),
        }
        //
        //it's easier to wait for server response instead of optimistic update here,
        //there are many edge cases here and this is the most reliable way
        MakeRequest(newBook, books, setBooks, toaster)

        // Reset form, might not be necessary
        setFormData({ title: '', author: '', publicationDate: '', totalPages: '' });
        setOpen(false)
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

            <button type="submit" className="bg-green-600 hover:bg-green-700 text-white rounded font-bold px-4 py-2">Add</button>
        </form>
    );
}

async function MakeRequest(Book, books, setBooks, toaster) {
    try {
        const response = await fetch('http://localhost:8046/books', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(Book)
        })

        if (response.ok) {
            const json = await response.json()
            setBooks((prevBooks) => [...prevBooks, json])
            toaster({
                title: 'Operation Successful',
                description: 'Book added successfully',
                variant: 'success',
            })
            return
        }
        else {
            setBooks(books)
            toaster({
                title: 'Operation Failed',
                description: 'Could not add Book',
                variant: 'destructive',
            })
        }
    }
    catch (error) {
        toaster({
            title: 'Operation Failed',
            description: 'Network Error, Could not connect to server',
            variant: 'destructive',
        })
        setBooks(books)
    }
}
