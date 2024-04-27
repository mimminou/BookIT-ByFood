"use client"
import { Button } from '@/components/ui/button'
import { Table, TableHeader, TableHead, TableRow, TableBody, TableCell } from '@/components/ui/table'
import AddBookDialog from '@/app/components/AddBookDialog'
import UpdateBookDialog from '@/app/components/UpdateBookDialog'
import DeleteBookDialog from '@/app/components/DeleteBookDialog'
import React, { useState, useEffect } from 'react'
import { useContext, Dispatch, SetStateAction } from 'react'
import { useRouter } from 'next/navigation'
import { Context } from '../context'
import { Book, ErrMessage } from '@/app/types'
import { toast } from '@/components/ui/use-toast'
import { MouseEvent } from 'react'


interface BookTableProps {
    books: Book[]
    setUpdateDialogOpen: Dispatch<SetStateAction<boolean>>
    setDeleteDialogOpen: Dispatch<SetStateAction<boolean>>
    setSelectedBook: Dispatch<SetStateAction<Book>>
}


export default function BookList() {

    const ctx = useContext(Context)
    const [error, setError] = useState(null)
    const [addDialogOpen, setAddDialogOpen] = useState(false)
    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
    const [updateDialogOpen, setUpdateDialogOpen] = useState(false)
    const [isLoading, setIsLoading] = useState(true)

    useEffect(() => {
        setIsLoading(true)
        GetBooks(ctx.toast).then((data: Book[] | ErrMessage) => {
            if (!data) {
                setError("Something went wrong")
                setIsLoading(false)
                return
            }
            if ("msg" in data) {
                setError(data.msg)
                setIsLoading(false)
                return
            }
            setIsLoading(false)
            ctx.setBooks(data)
        })
    }, [])

    return (
        <div className='p-4'>
            <h1 className='text-3xl text-gray-500'>Book List</h1>

            <UpdateBookDialog open={updateDialogOpen} setOpen={setUpdateDialogOpen} />
            <AddBookDialog open={addDialogOpen} setOpen={setAddDialogOpen} />
            <DeleteBookDialog open={deleteDialogOpen} setOpen={setDeleteDialogOpen} shouldRoute={false} />
            {isLoading ? null :
                <div>{
                    <Button className="bg-green-600 hover:bg-green-700" onClick={() => setAddDialogOpen(true)}>Add Book</Button>
                }</div>}
            {isLoading ?
                <p>Loading...</p> :
                <div className='flex items-center justify-center'>
                    <BookTable books={ctx.books} setUpdateDialogOpen={setUpdateDialogOpen} setDeleteDialogOpen={setDeleteDialogOpen} setSelectedBook={ctx.setSelectedBook}></BookTable>
                </div>
            }
        </div>
    )
}

function BookTable({ books, setUpdateDialogOpen,
    setDeleteDialogOpen, setSelectedBook }: BookTableProps) {

    const router = useRouter()
    const onUpdateClick = (event: MouseEvent<HTMLButtonElement>, book: Book) => {
        event.stopPropagation()
        setSelectedBook(book)
        setUpdateDialogOpen(true)
    }

    const onDeleteClick = (event: MouseEvent<HTMLButtonElement>, book: Book) => {
        event.stopPropagation()
        setSelectedBook(book)
        setDeleteDialogOpen(true)
    }

    return (
        <Table>
            <TableHeader>
                <TableRow>

                    <TableHead className="text-center">ID</TableHead>
                    <TableHead className="text-center">Title</TableHead>
                    <TableHead className="text-center">Author</TableHead>
                    <TableHead className="text-center">Publication Date</TableHead>
                    <TableHead className="text-center">Actions</TableHead>
                </TableRow>
            </TableHeader>
            <TableBody>

                {books.map((book: Book) => (
                    <TableRow key={book.book_id} className="cursor-pointer" onClick={() => {
                        router.push(`/books/${book.book_id}`)
                    }
                    } >
                        <TableCell className="text-center">{book.book_id} </TableCell>
                        <TableCell className="text-center">{book.title} </TableCell>
                        <TableCell className="text-center">{book.author}</TableCell>
                        <TableCell className="text-center">{book.pub_date}</TableCell>
                        <TableCell className="text-center flex justify-center gap-2 items-center">
                            <Button className="bg-blue-500 hover:bg-blue-700" onClick={(e) => onUpdateClick(e, book)}>Update</Button>
                            <Button className="bg-red-500 hover:bg-red-700" onClick={(e) => onDeleteClick(e, book)}>Delete</Button>
                        </TableCell>
                    </TableRow>
                ))}
            </TableBody>
        </Table>
    )
}

async function GetBooks(toaster: typeof toast) {
    try {
        const resp = await fetch("http://localhost:8046/books")
        let jsonResponse = await resp.json()
        if (!resp.ok) {
            //check if it's an err
            if ("msg" in jsonResponse) {
                jsonResponse = jsonResponse as ErrMessage
                throw new Error(jsonResponse.msg)
            }
            else {
                throw new Error("Unkown response type")
            }
        }
        jsonResponse = jsonResponse as Book[]
        jsonResponse.forEach((element: Book) => {
            element.pub_date = new Date(element.pub_date).toISOString().split('T')[0]
        })
        return jsonResponse
    }
    catch (error) {
        if ("msg" in error) {
            toaster({
                title: "Error, Could not fetch book list",
                description: error.msg,
                variant: "destructive",
            })
            return { msg: error.msg } as ErrMessage
        }
        toaster({
            title: "Error, Could not fetch book list",
            description: "Network Error",
            variant: "destructive",
        })
        return { msg: "Network Error" } as ErrMessage
    }
}
