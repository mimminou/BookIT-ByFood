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
                        <AddBookForm books={ctx.books} setBooks={ctx.setBooks} setOpen={setOpen} />
                    </DialogTitle>
                </DialogHeader>
            </DialogContent>
        </Dialog>
    )
}

function AddBookForm({ books, setBooks, setOpen }) {

    const [formData, setFormData] = useState({
        title: '',
        author: '',
        publicationDate: '',
        totalPages: '',
    });
    console.log(formData)

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

        // Submit form data (replace with your submission logic)
        console.log('Submitting form data:', formData);

        //Optimistic update
        const newBook = {
            title: formData.title,
            author: formData.author,
            pub_date: formData.publicationDate,
            num_pages: parseInt(formData.totalPages),
        }
        MakeRequest(newBook)

        // Reset form
        setFormData({ title: '', author: '', publicationDate: '', totalPages: '' });
        setOpen(false)
        const currentBooks = books
        setBooks([...currentBooks, newBook])
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

            <button type="submit">Add</button>
        </form>
    );
}

function MakeRequest(Book) {
    fetch('http://localhost:8046/books', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(Book)
    })
}
