// we define a book here because it's needed in all components in this project
interface Book {
    book_id?: number
    title: string
    author: string
    pub_date: string
    num_pages: number | string
}

interface ErrMessage {
    msg: string
}

class ServerError extends Error {
    constructor(Err : ErrMessage) {
        super(Err.msg);
        this.name = 'ServerError';
        this.msg = Err.msg
    }
    msg: string;
}

export type { Book, ErrMessage}
export {ServerError}
