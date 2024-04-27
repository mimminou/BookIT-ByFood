"use client"
import { Button } from '@/components/ui/button'
import BookCard from '@/app/components/BookCard'
import { Context } from '@/app/context'
import { useContext } from 'react'
import UpdateBookDialog from '@/app/components/UpdateBookDialog'
import DeleteBookDialog from '@/app/components/DeleteBookDialog'
import { useState, useEffect, Dispatch, SetStateAction } from 'react'
import { Book, ErrMessage, ServerError } from '@/app/types'
import { toast } from '@/components/ui/use-toast'
import { MouseEvent } from 'react'

export default function BookDetails({ params }) {

    const ctx = useContext(Context)
    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
    const [updateDialogOpen, setUpdateDialogOpen] = useState(false)
    const [book, setBook] = useState({} as Book)
    const [error, setError] = useState({ msg: "" } as ErrMessage)

    useEffect(() => {
        //maybe use a cached value here ?
        //I opted to make an API call every time, to at least use the getBook endpoint
        GetBook(params.bookID, ctx.toast).then((data: Book | ErrMessage) => {
            if ("msg" in data) {
                setError(data)
                return
            }
            setBook(data)
        })
    }, [])

    const UpdateButtonClick = (e: MouseEvent<HTMLButtonElement>) => {
        e.stopPropagation()
        // we need to use setSelected book so that we can reuse the UpdateBookDialog component
        ctx.setSelectedBook(book)
        setUpdateDialogOpen(true)
    }

    const DeleteButtonClick = (e: MouseEvent<HTMLButtonElement>) => {
        e.stopPropagation()
        // we need to use setSelected book so that we can reuse the UpdateBookDialog component
        ctx.setSelectedBook(book)
        setDeleteDialogOpen(true)
    }

    return (
        <div>
            <div className='p-4 m-4'>
                <BookCard book={book} />
                <Button className="bg-blue-500 hover:bg-blue-700 m-2" onClick={(e) => UpdateButtonClick(e)
                }> Update </Button>
                <Button className="bg-red-500 hover:bg-red-700 m-2" onClick={(e) => DeleteButtonClick(e)}> Delete </Button>
            </div>
            < UpdateBookDialog open={updateDialogOpen} setOpen={setUpdateDialogOpen} setBook={setBook} />
            <DeleteBookDialog open={deleteDialogOpen} setOpen={setDeleteDialogOpen} shouldRoute={true} />

        </div>
    );
}

async function GetBook(bookID: number, toaster: typeof toast) {
    try {
        const resp = await fetch(`http://localhost:8046/books/${bookID}`)
        let jsonResponse = await resp.json()
        if (!resp.ok) {
            //check if it's an err
            if ("msg" in jsonResponse) {
                jsonResponse = jsonResponse as ErrMessage
                throw new ServerError(jsonResponse)
            }
            else {
                throw new Error("Unkown response type")
            }
        }
        jsonResponse = jsonResponse as Book
        jsonResponse.pub_date = new Date(jsonResponse.pub_date).toISOString().split('T')[0]
        return jsonResponse
    }

    catch (error) {
        if ("msg" in error) {
            toaster({
                title: "Operation Failed",
                description: error.msg,
                variant: "destructive",
            })
            return error
        }
        toaster({
            title: "Operation Failed",
            description: "Network Error, Could not connect to server",
            variant: "destructive",
        })
        return { msg: "Network Error" } as ErrMessage
    }
}
