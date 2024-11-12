package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"net/http"
)

type Book struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Pages int `json:"pages"`
}

var bookshelf = []Book{
	{ID: 1, Name: "Blue Bird", Pages: 500},
}
var id_counter = 1

func getBooks(c *gin.Context) {
	c.JSON(200, bookshelf)
}

func getBook(c *gin.Context) {
	id := c.Param("id")
	for _, book := range bookshelf {
		if strconv.Itoa(book.ID) == id {
			c.JSON(200, book)
			return
		}
	}
	c.JSON(404, gin.H{"message": "book not found"})
}

func addBook(c *gin.Context) {
	var newBook Book
	if err := c.BindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	for _, book := range bookshelf {
		if book.Name == newBook.Name {
			c.JSON(http.StatusConflict, gin.H{"message": "duplicate book name"})
			return
		}
	}
	id_counter += 1
	newBook.ID = id_counter
	bookshelf = append(bookshelf, newBook)
	c.JSON(http.StatusCreated, newBook)
}

func deleteBook(c *gin.Context) {
	id := c.Param("id")
	for i, book := range bookshelf {
		if strconv.Itoa(book.ID) == id {
			bookshelf = append(bookshelf[:i], bookshelf[i+1:]...)
			c.JSON(http.StatusNoContent, nil)
			return
		}
	}
	c.JSON(204, gin.H{"message": "book not found"})
}

func updateBook(c *gin.Context) {
	id := c.Param("id")
	var updatedBook Book
	if err := c.BindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	for _, book := range bookshelf {
		if book.Name == updatedBook.Name {
			c.JSON(http.StatusConflict, gin.H{"message": "duplicate book name"})
			return
		}
	}

	for i, book := range bookshelf {
		if strconv.Itoa(book.ID) == id {
			bookshelf[i].Name = updatedBook.Name
			bookshelf[i].Pages = updatedBook.Pages
			c.JSON(http.StatusOK, bookshelf[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "book not found"})
}

func main() {
	r := gin.Default()
	r.RedirectFixedPath = true

	// TODO: Add routes
	r.GET("/bookshelf", getBooks)
	r.GET("/bookshelf/:id", getBook)
	r.POST("/bookshelf", addBook)
	r.DELETE("/bookshelf/:id", deleteBook)
	r.PUT("/bookshelf/:id", updateBook)

	err := r.Run(":8087")
	if err != nil {
		return
	}
}
