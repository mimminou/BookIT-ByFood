"use client"

import { Button } from '@/components/ui/button'
import BookCard from '../../components/BookCard'
import { useState, useEffect } from 'react'

export default function BookDetails({ params }) {

    const [book, setBook] = useState({})
    const [error, setError] = useState(null)

    useEffect(() => {
        GetBook(params.bookID).then((data) => {
            if (!data) {
                setError("Something went wrong")
                return
            }
            setBook(data)
        })
    }, [])

    const UpdateButtonClick = (e) => {
        e.stopPropagation()
        alert("Update")
    }

    const DeleteButtonClick = (e) => {
        e.stopPropagation()
        alert("Delete")
    }


    return (
        <div className='p-4 m-4'>
            <BookCard book={book} />
            <Button className="bg-blue-500 hover:bg-blue-700 m-2" onClick={UpdateButtonClick}>Update</Button>
            <Button className="bg-red-500 hover:bg-red-700 m-2" onClick={DeleteButtonClick}>Delete</Button>
        </div>
    );
}

async function GetBook(bookID) {
    let jsonResponse = {}
    try {
        const resp = await fetch(`http://localhost:8046/books/${bookID}`)
        jsonResponse = await resp.json()
        jsonResponse.pub_date = new Date(jsonResponse.pub_date).toDateString()
    }
    catch {
        //show an err here
        console.log("err")
    }
    finally {
        return jsonResponse
    }
}
