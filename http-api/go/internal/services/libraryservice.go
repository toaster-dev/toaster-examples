package services

import (
	"context"
	"errors"
	"fmt"

	"toasterexample/internal/stores"
	"toasterexample/internal/xerrors"
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

func (s *LibraryService) GetBook(ctx context.Context, bookID uuid.UUID) (types.Book, error) {
	bookEntity, err := s.bookStore.GetBook(ctx, bookID)
	if errors.Is(err, stores.ErrNotFound) {
		return types.Book{}, xerrors.WrapError(xerrors.ErrNotFound, "book not found", err)
	} else if err != nil {
		return types.Book{}, fmt.Errorf("failed to get book: %w", err)
	}

	return types.BookEntityToBook(bookEntity), nil
}

func (s *LibraryService) CreateBook(ctx context.Context, book types.Book) (types.Book, error) {
	bookEntity, err := s.bookStore.CreateBook(ctx, types.BookToBookEntity(book))
	if err != nil {
		return types.Book{}, fmt.Errorf("failed to create book: %w", err)
	}

	return types.BookEntityToBook(bookEntity), nil
}
