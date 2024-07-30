package books

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Frankkwizera/go-gin-api-medium/pkg/common/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Define a mock database handler
type MockDB struct {
	mock.Mock
	*gorm.DB
}

func TestAddBook(t *testing.T) {
	// Set up Gin router
	r := gin.Default()

	// Use SQLite in-memory database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to the database: %v", err)
	}

	// Auto migrate Book model
	db.AutoMigrate(&models.Book{})

	// Create handler
	h := &handler{DB: db}

	// Register routes
	r.POST("/books", h.AddBook)

	// Create a sample book request
	newBook := AddBookRequestBody{
		Title:       "Test Title",
		Author:      "Test Author",
		Description: "Test Description",
	}

	// Convert book to JSON
	body, _ := json.Marshal(newBook)

	// Create a request to pass to our handler
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check the status code
	assert.Equal(t, http.StatusCreated, w.Code)

	// Check the response body
	var createdBook models.Book
	err = json.Unmarshal(w.Body.Bytes(), &createdBook)
	assert.Nil(t, err)
	assert.Equal(t, newBook.Title, createdBook.Title)
	assert.Equal(t, newBook.Author, createdBook.Author)
	assert.Equal(t, newBook.Description, createdBook.Description)

	// Verify the book was saved in the database
	var dbBook models.Book
	if err := db.First(&dbBook, createdBook.ID).Error; err != nil {
		t.Fatalf("book not found in database: %v", err)
	}
	assert.Equal(t, newBook.Title, dbBook.Title)
	assert.Equal(t, newBook.Author, dbBook.Author)
	assert.Equal(t, newBook.Description, dbBook.Description)
}

func TestDeleteBook(t *testing.T) {
	// Set up Gin router
	r := gin.Default()

	// Use SQLite in-memory database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to the database: %v", err)
	}

	// Auto migrate Book model
	db.AutoMigrate(&models.Book{})

	// Create a sample book in the database
	book := models.Book{
		Title:       "Test Book",
		Author:      "Test Author",
		Description: "Test Description",
	}
	db.Create(&book)

	// Create handler
	h := &handler{DB: db}

	// Register routes
	r.DELETE("/books/:id", h.DeleteBook)

	// Create a request to delete the book
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/books/%d", book.ID), nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Check the database to ensure the book is deleted
	var deletedBook models.Book
	result := db.First(&deletedBook, book.ID)
	assert.Error(t, result.Error)
	assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
}

func TestGetBook(t *testing.T) {
	// Set up Gin router
	r := gin.Default()

	// Use SQLite in-memory database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to the database: %v", err)
	}

	// Auto migrate Book model
	db.AutoMigrate(&models.Book{})

	// Create a sample book in the database
	book := models.Book{
		Title:       "Test Book",
		Author:      "Test Author",
		Description: "Test Description",
	}
	db.Create(&book)

	// Create handler
	h := &handler{DB: db}

	// Register routes
	r.GET("/books/:id", h.GetBook)

	// Create a request to get the book
	req, _ := http.NewRequest("GET", fmt.Sprintf("/books/%d", book.ID), nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Check the response body
	var fetchedBook models.Book
	err = json.Unmarshal(w.Body.Bytes(), &fetchedBook)
	assert.Nil(t, err)
	assert.Equal(t, book.ID, fetchedBook.ID)
	assert.Equal(t, book.Title, fetchedBook.Title)
	assert.Equal(t, book.Author, fetchedBook.Author)
	assert.Equal(t, book.Description, fetchedBook.Description)
}

func TestGetBooks(t *testing.T) {
	// Set up Gin router
	r := gin.Default()

	// Use SQLite in-memory database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to the database: %v", err)
	}

	// Auto migrate Book model
	db.AutoMigrate(&models.Book{})

	// Create some sample books in the database
	books := []models.Book{
		{Title: "Book 1", Author: "Author 1", Description: "Description 1"},
		{Title: "Book 2", Author: "Author 2", Description: "Description 2"},
	}
	db.Create(&books)

	// Create handler
	h := &handler{DB: db}

	// Register routes
	r.GET("/books", h.GetBooks)

	// Create a request to get the list of books
	req, _ := http.NewRequest("GET", "/books", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Check the response body
	var fetchedBooks []models.Book
	err = json.Unmarshal(w.Body.Bytes(), &fetchedBooks)
	assert.Nil(t, err)
	assert.Equal(t, len(books), len(fetchedBooks))

	// Verify that each book matches the expected data
	for i, book := range books {
		assert.Equal(t, book.Title, fetchedBooks[i].Title)
		assert.Equal(t, book.Author, fetchedBooks[i].Author)
		assert.Equal(t, book.Description, fetchedBooks[i].Description)
	}
}

func TestUpdateBook(t *testing.T) {
	// Set up Gin router
	r := gin.Default()

	// Use SQLite in-memory database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to the database: %v", err)
	}

	// Auto migrate Book model
	db.AutoMigrate(&models.Book{})

	// Create a sample book in the database
	book := models.Book{
		Title:       "Original Title",
		Author:      "Original Author",
		Description: "Original Description",
	}
	db.Create(&book)

	// Create handler
	h := &handler{DB: db}

	// Register routes
	r.PUT("/books/:id", h.UpdateBook)

	// Create an updated book request
	updatedBook := UpdateBookRequestBody{
		Title:       "Updated Title",
		Author:      "Updated Author",
		Description: "Updated Description",
	}
	body, _ := json.Marshal(updatedBook)

	// Create a request to update the book
	// req, _ := http.NewRequest("PUT", "/books/"+string(book.ID), bytes.NewBuffer(body))
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/books/%d", book.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Check the response body
	var fetchedBook models.Book
	err = json.Unmarshal(w.Body.Bytes(), &fetchedBook)
	assert.Nil(t, err)
	assert.Equal(t, book.ID, fetchedBook.ID)
	assert.Equal(t, updatedBook.Title, fetchedBook.Title)
	assert.Equal(t, updatedBook.Author, fetchedBook.Author)
	assert.Equal(t, updatedBook.Description, fetchedBook.Description)

	// Verify that the book was updated in the database
	var dbBook models.Book
	if err := db.First(&dbBook, book.ID).Error; err != nil {
		t.Fatalf("book not found in database: %v", err)
	}
	assert.Equal(t, updatedBook.Title, dbBook.Title)
	assert.Equal(t, updatedBook.Author, dbBook.Author)
	assert.Equal(t, updatedBook.Description, dbBook.Description)
}
