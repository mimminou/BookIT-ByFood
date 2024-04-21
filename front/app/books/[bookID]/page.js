"use client"

import { Button } from '@/components/ui/button'
import BookCard from '@/app/components/BookCard'
import { Context } from '@/app/context'
import { useContext } from 'react'
import UpdateBookDialog from '@/app/components/UpdateBookDialog'
import DeleteBookDialog from '@/app/components/DeleteBookDialog'
import { useState, useEffect } from 'react'

export default function BookDetails({ params }) {

    const ctx = useContext(Context)
    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
    const [updateDialogOpen, setUpdateDialogOpen] = useState(false)
    const [book, setBook] = useState({})
    const [error, setError] = useState(null)

    useEffect(() => {
        //maybe use a cached value here ?
        //I opted to make an API call every time, to at least use the getBook endpoint
        GetBook(params.bookID, ctx.toast).then((data) => {
            if (!data) {
                setError("Something went wrong")
                return
            }
            setBook(data)
        })
    }, [])

    const UpdateButtonClick = (e) => {
        e.stopPropagation()
        // we need to use setSelected book so that we can reuse the UpdateBookDialog component
        ctx.setSelectedBook(book)
        setUpdateDialogOpen(true)
    }

    const DeleteButtonClick = (e) => {
        e.stopPropagation()
        // we need to use setSelected book so that we can reuse the UpdateBookDialog component
        ctx.setSelectedBook(book)
        setDeleteDialogOpen(true)
    }

    return (
        <div>
            <div className='p-4 m-4'>
                <BookCard book={book} />
                <Button className="bg-blue-500 hover:bg-blue-700 m-2" onClick={(e) => UpdateButtonClick(e)}>Update</Button>
                <Button className="bg-red-500 hover:bg-red-700 m-2" onClick={(e) => DeleteButtonClick(e)}>Delete</Button>
            </div>
            <UpdateBookDialog open={updateDialogOpen} setOpen={setUpdateDialogOpen} setBook={setBook} />
            <DeleteBookDialog open={deleteDialogOpen} setOpen={setDeleteDialogOpen} shouldRoute={true} />

        </div>
    );
}

async function GetBook(bookID, toaster) {
    let jsonResponse = {}
    try {
        const resp = await fetch(`http://localhost:8046/books/${bookID}`)
        if (!resp.ok) {
            toaster({
                title: "Operation Failed",
                description: ", Could not fetch book",
                variant: "destructive",
            })
            return null
        }
        jsonResponse = await resp.json()
        jsonResponse.pub_date = new Date(jsonResponse.pub_date).toISOString().split('T')[0]
    }

    catch {
        //show an err here
    }
    finally {
        return jsonResponse
    }
}
