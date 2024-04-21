## Golang Server for the Book Library assignment
# Running
## First time running
cd into the backend folder (/back/) and run `go run . -s` to setup a local sqlite DB with the schema required for the server to function.

once the setup process is finished, you can continue to the next step

## Running the server
use `go run .` to run the server from terminal.

The server has an stdout to the console that prints incoming requests and their responses with Timestamp, Endpoint, HTTP Method as well as the body of the request if available.

## Project Structure

- `/DB/*`: contains the sqlite database and it's schema
- `/server/*`: contains the server, controllers and middleware implementations
- `/utils/*`: contains utility functions
- `/config.json`: configuration file for the server, specifies the port and the path to the DB and it's schema


## API Endpoints

### Models:
- Book Json Structure : 

`{"book_id": int, "title": string, "author": string, "num_pages": int, "pub_date": Date string*"YYYY-MM-DD"*}`

### Endpoints:
- `/books`: `GET` get all books, returns a json array of Book objects or an error message
- `/books/`: `POST` create a new book, takes in a json object of type Book (without book_id key) and returns the created book as json or an error message
- `/books/{id}`: `GET` get a specific book by id, returns a json object of Book or an error message if not found
- `/books/{id}`: `PUT` update a specific book by id, takes in a json object of type Book and returns the updated book as json or an error message if not found
- `/books/{id}`: `DELETE` delete a specific book by id, returns a success message or an error message if not found

# Tests

Unit tests were written for both /server and /utils.

simply navigate to any of these folders and run `go test -v`

This will run all the tests available on that directory:

- server and controllers for server/
- utility functions for utils/
