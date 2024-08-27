package types

import (
	"time"

	"toasterexample/internal/stores/entities"

	"github.com/google/uuid"
)

type Book struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func BookEntityToBook(e entities.Book) Book {
	return Book(e)
}
