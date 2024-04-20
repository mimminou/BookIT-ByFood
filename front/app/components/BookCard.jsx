import { Card, CardContent, CardTitle, CardHeader } from "@/components/ui/card"
import Image from 'next/image'

export default function BookCard({ book }) {
    return (
        <Card className="min-w-1/6">
            <CardHeader>
                <CardTitle>{book.title}</CardTitle>
            </CardHeader>
            <CardContent>
                <BookDetails book={book} />
            </CardContent>
        </Card>
    )
}

function BookDetails({ book }) {
    return (
        <div className="flex flex-col md:flex-row justify-between">
            <ul>
                <li>Book ID : {book.book_id} </li>
                <li>Book author : {book.author} </li>
                <li>Book publication date : {book.pub_date} </li>
                <li>Book number of pages : {book.num_pages ? book.num_pages : "Not specified"} </li>
            </ul>
            <Image src="/ByFoodCover.png" width={200} height={200} alt="ByFood" />
        </div>
    )
}
