package stores

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"toasterexample/internal/stores/entities"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/sqlscan"
	"github.com/google/uuid"
)

type BookStore struct {
	db *sql.DB
}

func NewBookStore(db *sql.DB) *BookStore {
	return &BookStore{
		db: db,
	}
}

func (s *BookStore) ListBooks(ctx context.Context, lastID uuid.UUID, limit int) ([]entities.Book, bool, error) {
	books := []entities.Book{}

	q := sq.Select("id", "title", "created_at", "updated_at").
		From("books").
		OrderBy("id DESC").
		Limit(uint64(limit + 1))

	if lastID != uuid.Nil {
		q = q.Where(sq.Lt{"id": lastID})
	}

	query, args, err := q.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return []entities.Book{}, false, fmt.Errorf("failed to build query: %w", err)
	}

	if err := sqlscan.Select(ctx, s.db, &books, query, args...); err != nil {
		return []entities.Book{}, false, fmt.Errorf("failed to list books: %w", err)
	}

	hasMore := false
	if len(books) > limit {
		hasMore = true
		books = books[:limit]
	}

	return books, hasMore, nil
}

func (s *BookStore) GetBook(ctx context.Context, bookID uuid.UUID) (entities.Book, error) {
	query := `
		SELECT id, title, created_at, updated_at
		FROM books
		WHERE id = $1
	`

	var book entities.Book
	if err := sqlscan.Get(ctx, s.db, &book, query, bookID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Book{}, ErrNotFound
		}

		return entities.Book{}, fmt.Errorf("failed to get book: %w", err)
	}

	return book, nil
}

func (s *BookStore) CreateBook(ctx context.Context, book entities.Book) (entities.Book, error) {
	bookID, err := uuid.NewV7()
	if err != nil {
		return entities.Book{}, fmt.Errorf("failed to generate book ID: %w", err)
	}

	book.ID = bookID
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()

	query := `
		INSERT INTO books (id, title, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`

	res, err := s.db.ExecContext(ctx, query, book.ID, book.Title, book.CreatedAt, book.UpdatedAt)
	if err != nil {
		return entities.Book{}, fmt.Errorf("failed to create book: %w", err)
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected != 1 {
		return entities.Book{}, fmt.Errorf("failed to create book: expected 1 row to be affected, got %d", rowsAffected)
	}

	return book, nil
}
