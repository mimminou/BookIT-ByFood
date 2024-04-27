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
- `/database/*`: contains the database interface implementation
- `/server/*`: contains the server and middlwares
- `/controllers/*`: contains the controllers for the endpoints
- `/models/*`: contains the models used in the project
- `/services/*`: contains the logic for each endpoint
- `/utils/*`: contains utility functions
- `/config.json`: configuration file for the server, specifies the port and the path to the DB and it's schema

## Books :
### Models:
- Book Json Structure : 

```
{"book_id": int, "title": string, "author": string, "num_pages": int, "pub_date": Date string*"YYYY-MM-DD"*}
```

### Endpoints:
- `/books`: `GET` get all books, returns a json array of Book objects or an error message
- `/books/`: `POST` create a new book, takes in a json object of type Book (without book_id key) and returns the created book as json or an error message
- `/books/{id}`: `GET` get a specific book by id, returns a json object of Book or an error message if not found
- `/books/{id}`: `PUT` update a specific book by id, takes in a json object of type Book and returns the updated book as json or an error message if not found
- `/books/{id}`: `DELETE` delete a specific book by id, returns a success message or an error message if not found


## Url Cleaner :
### Models:

- Request Json Structure : 

```
{
"url" : String (valid URL),
"operation" : string (one of the following : "canonical", "redirection" "all")
}

```

### Response Json Structure :
- on success :   `{"processed_url" : string}`
- on error :  `{"msg": string}`


### Endpoints:
- `/url`: `POST` : takes a JSON body with the Request structure (as specified above)

### Operations:

##### `canonical` : Cleans up URL, removes Query parameters and trailing slashes.

example : `https://BYFOOD.com/food-EXPeriences?query=abc/`  -->  `https://BYFOOD.com/food-EXPeriences`



##### `redirection` : Checks if the domain is `byfood.com` (as a security measure, will only work for www subdomain, and add it to the URL if it's missing), then returns the url lowercased.

example : `https://BYFOOD.com/food-EXPeriences?query=abc/` --> `https://www.byfood.com/food-experiences?query=abc/`


##### `all` : Performs both Canonical and Redirection, as such, it's bound to only work on `byfood.com` domain

example : `https://BYFOOD.com/food-EXPeriences?query=abc/` --> `https://www.byfood.com/food-experiences`



### Interactive docs
Interactive docs were generated using Swagger for this project, simply run the server and access the root address (http://localhost:PORT/)


# Tests

Unit tests were written for all components of the backend

simply navigate to any of these folders and run `go test -v`

This will run all the tests available on that directory

To run all tests, simply run `go test -v ./...` from  `/back` directory
