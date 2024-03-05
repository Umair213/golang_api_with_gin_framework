package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

func getBooks(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, books)
}

func deleteBookByID(context *gin.Context) {
	id := context.Param("id")
	book, err := delBookByID(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found!"})
	}
	context.IndentedJSON(http.StatusNoContent, book)
}

func delBookByID(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			deletedBook := &books[i]
			fmt.Println(deletedBook)
			books = append(books[:i], books[i+1:]...)
			return deletedBook, nil
		}
	}
	return nil, errors.New("book not found")
}

func checkoutBook(context *gin.Context) {
	id, okay := context.GetQuery("id")
	if !okay {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter!"})
		return
	}
	book, err := getBookByID(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found!"})
	}

	if book.Quantity <= 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available."})
		return
	}

	book.Quantity -= 1
	context.IndentedJSON(http.StatusOK, book)
}

func returnBook(context *gin.Context) {
	id := context.Param("id")
	book, err := getBookByID(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found!"})
	}

	book.Quantity += 1
	context.IndentedJSON(http.StatusOK, book)
}

func bookByID(context *gin.Context) {
	id := context.Param("id")
	book, err := getBookByID(id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found!"})
	}
	context.IndentedJSON(http.StatusOK, book)
}

func getBookByID(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("book not found")
}

func createBook(context *gin.Context) {
	var newBook book
	if err := context.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	context.IndentedJSON(http.StatusCreated, books)
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.POST("/create_book", createBook)
	router.GET("/book/:id", bookByID) //Setting path parameter id
	router.DELETE("/book/:id", deleteBookByID)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return_book/:id", returnBook)
	router.Run("localhost:8080")
}
