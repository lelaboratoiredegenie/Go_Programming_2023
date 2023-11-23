package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Book struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Pages int    `json:"pages"`
}

var bookshelf = []Book{
	{1, "Blue Bird", 500},
}

var count int = 1

func getBooks(c *gin.Context) {
	c.JSON(http.StatusOK, bookshelf)

}
func getBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	fmt.Println(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	if id <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	var foundBook Book
	for _, book := range bookshelf {
		if book.ID == id {
			foundBook = book
			break
		}
	}

	if foundBook.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "book not found"})
		return
	}

	c.JSON(http.StatusOK, foundBook)
}
func addBook(c *gin.Context) {
	var newBook Book

	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the maximum ID in the existing bookshelf
	for _, book := range bookshelf {
		if book.Name == newBook.Name {
			c.JSON(http.StatusConflict, gin.H{"message": "duplicate book name"})
			return
		}
	}
	newBook.ID = count + 1

	count = count + 1

	bookshelf = append(bookshelf, newBook)
	c.JSON(http.StatusCreated, newBook)
}

func deleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	if id <= 0 {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	var foundIndex = -1
	for i, book := range bookshelf {
		if book.ID == id {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	bookshelf = append(bookshelf[:foundIndex], bookshelf[foundIndex+1:]...)
	c.JSON(http.StatusNoContent, gin.H{})
}

func updateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	if id <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "book not found"})
		return
	}

	var updatedBook Book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var foundIndex = -1
	for i, book := range bookshelf {
		if book.Name == updatedBook.Name {
			c.JSON(http.StatusConflict, gin.H{"message": "duplicate book name"})
			return
		}
		if book.ID == id {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		c.JSON(http.StatusNotFound, gin.H{"message": "book not found"})
		return
	}

	updatedBook.ID = id
	bookshelf[foundIndex] = updatedBook
	c.JSON(http.StatusOK, updatedBook)
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
