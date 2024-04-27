import { Dialog, DialogHeader, DialogContent, DialogTitle } from "@/components/ui/dialog"
import { Context } from '../context'
import { Dispatch, SetStateAction, useContext, FormEvent } from 'react'
import { useRouter } from 'next/navigation'
import { Book, ErrMessage } from '@/app/types'
import { toast } from '@/components/ui/use-toast'


interface DeleteBookFormProps {
    books: Book[]
    setBooks: Dispatch<SetStateAction<Book[]>>
    setOpen: Dispatch<SetStateAction<boolean>>
    book: Book
    toaster: typeof toast
    shouldRoute: boolean
}

interface RequestArgs {
    book: Book
    books: Book[]
    setBooks: Dispatch<SetStateAction<Book[]>>
    toaster: typeof toast
    shouldRoute: boolean
    route: Function
}

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


function DeleteBookForm(props: DeleteBookFormProps) {
    const { books, setBooks, setOpen, book, toaster, shouldRoute } = props
    const router = useRouter()
    const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();

        MakeRequest({
            book,
            books,
            setBooks,
            toaster,
            shouldRoute,
            route: () => {
                console.log('Book deleted')
                router.push('/books') // Need to do it this way, as it's not possible to pass the router to non components
            }
        }
        )

        const updatedBookList = books.filter((element: Book) => {
            if (element.book_id !== book.book_id) {
                return element
            }
        })
        setOpen(false)
        setBooks(updatedBookList)
    }

    return (
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">

            <p className="font-bold">{book.title}</p>

            <p>Are you sure you want to delete this book?</p>

            <button type="submit" className="bg-red-500 hover:bg-red-700 text-white rounded font-bold px-4 py-2">Delete</button>
        </form>
    );
}

async function MakeRequest({ book, books, setBooks, toaster, shouldRoute, route }: RequestArgs) {
    const oldBooks = books // save old state
    try {
        const response = await fetch(`http://localhost:8046/books/${book.book_id}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
            },
        })

        if (!response.ok) {
            let jsonResponse = await response.json()
            console.log("RESPNSOE IS NOT OK", response.status)
            //check if it's an err
            if ("msg" in jsonResponse) {
                jsonResponse = jsonResponse as ErrMessage
                throw new Error(jsonResponse.msg)
            }
            else {
                throw new Error("Unkown response type")
            }
        }
        toaster({
            description: 'Book Deleted',
        })
        if (shouldRoute) {
            route()
        }
    }
    catch (error) {
        setBooks(oldBooks)
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
