import type { Book } from '@/app/types'
export default function BookComponent({ book }: { book: Book }) {
    return (
        <div>
            <ul>
                <li>Book title : {book.title} </li>
                <li>Book author : {book.author} </li>
                <li>Book number of pages : {book.num_pages ? book.num_pages : "Not specified"} </li>
                <li>Book publication date : {book.pub_date} </li>
            </ul>
        </div>
    )
}
