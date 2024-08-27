package services

import (
	"context"
	"errors"
	"fmt"

	"toasterexample/internal/stores"
	"toasterexample/types"

	"github.com/google/uuid"
)

type LibraryService struct {
	bookStore *stores.BookStore
}

func NewLibraryService(bookStore *stores.BookStore) *LibraryService {
	return &LibraryService{
		bookStore: bookStore,
	}
}

func (s *LibraryService) ListBooks(ctx context.Context, lastID uuid.UUID, limit int) ([]types.Book, bool, error) {
	booksEntities, hasMore, err := s.bookStore.ListBooks(ctx, lastID, limit)
	if err != nil {
		return []types.Book{}, false, fmt.Errorf("failed to list books: %w", err)
	}

	books := make([]types.Book, 0, len(booksEntities))
	for _, bookEntity := range booksEntities {
		books = append(books, types.BookEntityToBook(bookEntity))
	}

	return books, hasMore, nil
}

