"use client"
import { createContext, useState } from 'react'

export const Context = createContext()

const CtxProvider = ({ children }) => {
    const [books, setBooks] = useState([])
    const [selectedBookID, setSelectedBookID] = useState(0)
    const [selectedBook, setSelectedBook] = useState({ book_id: 0, title: "", author: "", pub_date: "", num_pages: "" })

    return (
        <Context.Provider value={{ books, setBooks, selectedBookID, setSelectedBookID, selectedBook, setSelectedBook }}>
            {children}
        </Context.Provider>
    )
}

export default CtxProvider
