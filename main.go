package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string `json:"id"` //it musn't have empty spaces at json:<"field">
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

// Hardcoded books, this also should come from database
var books = []book{
	{ID: "1", Title: "In search of lost time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "In search of lost time 2", Author: "Marcel Proust", Quantity: 5},
	{ID: "3", Title: "In search of lost time 3", Author: "Marcel Proust", Quantity: 3},
}

// getBook accepts gin context as parameter. What is the *??. Context is like "request"
func getBooks(context *gin.Context) {

	//With IndentedJSON basically we are returning book in nice JSON format. It's like jsonify from flask.
	context.IndentedJSON(http.StatusOK, books)
}

func createBook(context *gin.Context) {
	var newBook book

	//nil = null.
	//& return the address of the variable
	//Bind json will verify the incoming data that is a book type and bind it to newBook address
	//If err is not null, means there is an error then return. BindJSON will return an error.
	if err := context.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook) //appends will add a book to the end of the slice

	context.IndentedJSON(http.StatusCreated, newBook)
}

// We will return a pointer to a book, and an error that could be nil (null)
func getBookByIdHelper(id string) (*book, error) {
	//Loop through book
	for index, book := range books {
		if book.ID == id {
			//If we find a book, then return the memory address of that book and nil at error
			return &books[index], nil
		}
	}
	return nil, errors.New("book not found")
}

func getBookById(context *gin.Context) {
	//Extract the id from param
	id := context.Param("id")
	//call the helper
	book, err := getBookByIdHelper(id)
	if err != nil {
		return
	}
	context.IndentedJSON(http.StatusOK, book)
}

func checkoutBook(context *gin.Context) {
	id, ok := context.GetQuery("id")
	if !ok {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad id from query params"})
		return
	}

	book, err := getBookByIdHelper(id)

	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not found"})
		return
	}

	if book.Quantity == 0 {
		//With sprintf we can combine string with variables
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Not enough stock of the book. Stock: %v", book.Quantity)})
		return
	}

	book.Quantity -= 1
	context.IndentedJSON(http.StatusOK, gin.H{"message": "Book checkout successfully", "book": book})
}

func main() {
	router := gin.Default()

	//Define the route and the function for that route
	router.GET("/books", getBooks)

	router.GET("/books/:id", getBookById)

	router.PATCH("/checkout/books", checkoutBook)

	router.POST("/books", createBook)

	router.Run("localhost:8080")
}
