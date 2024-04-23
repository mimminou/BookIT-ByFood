import { Dialog, DialogHeader, DialogContent, DialogTitle } from "@/components/ui/dialog"
import { Context } from '../context'
import { useContext } from 'react'
import { useRouter } from 'next/navigation'


export default function DeleteBookDialog({ open, setOpen, shouldRoute }) {
    const ctx = useContext(Context)
    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>
                        Delete Book
                    </DialogTitle>
                </DialogHeader>
                <DeleteBookForm books={ctx.books} setBooks={ctx.setBooks}
                    setOpen={setOpen} book={ctx.selectedBook}
                    toaster={ctx.toast} shouldRoute={shouldRoute} />
            </DialogContent>
        </Dialog>
    )
}

function DeleteBookForm({ books, setBooks, setOpen, book, toaster, shouldRoute }) {

    const router = useRouter()

    const handleSubmit = (event) => {
        event.preventDefault();

        MakeRequest(book, books, setBooks, toaster, router, shouldRoute)

        const updatedBookList = books.filter((element) => {
            if (element.book_id !== book.book_id) {
                return element
            }
        })
        setOpen(false)
        setBooks(updatedBookList)
    };

    return (
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">

            <p className="font-bold">{book.title}</p>

            <p>Are you sure you want to delete this book?</p>

            <button type="submit" className="bg-red-500 hover:bg-red-700 text-white rounded font-bold px-4 py-2">Delete</button>
        </form>
    );
}

async function MakeRequest(book, books, setBooks, toaster, router, shouldRoute) {
    const oldBooks = books // save old state
    try {
        const response = await fetch(`http://localhost:8046/books/${book.book_id}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
            },
        })

        if (response.ok) {
            toaster({
                description: 'Book Deleted',
            })

            if (shouldRoute) {
                router.push('/books')
            }
            return
        }
        else {
            const json = await response.json()
            setBooks(oldBooks)
            toaster({
                title: 'Operation Failed',
                description: json.msg ? json.msg : 'Deleting Book failed',
                variant: 'destructive',
            })
        }
    }
    catch (error) {
        setBooks(oldBooks)
        toaster({
            title: 'Operation Failed',
            description: 'Network Error, Could not connect to server',
            variant: 'destructive',
        })
        return
    }
}
