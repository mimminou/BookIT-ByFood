"use client"
import { createContext, useState, Dispatch, SetStateAction } from 'react'
import { Toaster } from '@/components/ui/toaster'
import { useToast, toast } from '@/components/ui/use-toast'
import { Book } from '@/app/types'


type Ctx = {
    books: Book[]
    setBooks: Dispatch<SetStateAction<Book[]>>
    selectedBookID: number
    setSelectedBookID: Dispatch<SetStateAction<number>>
    selectedBook: Book
    setSelectedBook: Dispatch<SetStateAction<Book>>
    toast: typeof toast
}

export const Context = createContext({} as Ctx)


const CtxProvider = ({ children }) => {
    const [books, setBooks] = useState([])
    const [selectedBookID, setSelectedBookID] = useState(0)
    const [selectedBook, setSelectedBook] = useState({ book_id: 0, title: "", author: "", pub_date: "", num_pages: "" })
    const { toast } = useToast()

    return (
        <Context.Provider value={{ books, setBooks, selectedBookID, setSelectedBookID, selectedBook, setSelectedBook, toast }}>
            {children}
            <Toaster />
        </Context.Provider>
    )
}

export default CtxProvider
