"use client"
import { createContext, useState } from 'react'
import { Toaster } from '@/components/ui/toaster'
import { useToast } from '@/components/ui/use-toast'

export const Context = createContext()

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
