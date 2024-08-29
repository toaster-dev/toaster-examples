package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"toasterexample/internal/services"
	"toasterexample/internal/xerrors"
	"toasterexample/types"

	"github.com/google/uuid"
)

func ListBooks(libraryService *services.LibraryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		lastID := uuid.Nil
		if lastIDStr := r.URL.Query().Get("lastID"); lastIDStr != "" {
			id, err := uuid.Parse(lastIDStr)
			if err != nil {
				slog.Error("invalid lastID", slog.Any("error", err), slog.String("lastID", lastIDStr))
				writeJSON(w, http.StatusBadRequest, httpError{
					Message: "invalid lastID, must be a valid UUID",
				})
				return
			}

			lastID = id
		}

		limit := 10
		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			l, err := strconv.Atoi(limitStr)
			if err != nil {
				slog.Error("invalid limit", slog.Any("error", err), slog.String("limit", limitStr))
				writeJSON(w, http.StatusBadRequest, httpError{
					Message: "invalid limit, must be an integer",
				})
				return
			}

			if l <= 0 {
				slog.Error("invalid limit", slog.String("limit", limitStr))
				writeJSON(w, http.StatusBadRequest, httpError{
					Message: "invalid limit, must be greater than 0",
				})
				return
			}

			limit = l
		}

		books, hasMore, err := libraryService.ListBooks(ctx, lastID, limit)
		if err != nil {
			slog.Error("failed to list books", slog.Any("error", err))
			writeJSON(w, http.StatusInternalServerError, httpError{
				Message: "failed to list books",
			})
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"data": books,
			"metadata": map[string]any{
				"pagination": types.PaginationCursor{
					HasMore: hasMore,
				},
			},
		})
	}
}

func GetBook(libraryService *services.LibraryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		bookID, err := uuid.Parse(r.PathValue("bookID"))
		if err != nil {
			slog.Error("invalid book ID", slog.Any("error", err), slog.String("bookID", r.PathValue("bookID")))
			writeJSON(w, http.StatusBadRequest, httpError{
				Message: "invalid book ID, must be a valid UUID",
			})
			return
		}

		book, err := libraryService.GetBook(ctx, bookID)
		if err != nil {
			slog.Error("failed to get book", slog.Any("error", err), slog.String("bookID", bookID.String()))

			status := http.StatusInternalServerError
			message := "failed to get book"

			switch {
			case errors.Is(err, xerrors.ErrNotFound):
				status = http.StatusNotFound

			default:
				writeJSON(w, http.StatusInternalServerError, httpError{
					Message: "failed to get book",
				})
			}

			var serr xerrors.StructuredError
			if errors.As(err, &serr) {
				message = err.Error()
			}

			writeJSON(w, status, httpError{
				Message: message,
			})
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"data": book,
		})
	}
}

func CreateBook(libraryService *services.LibraryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var requestCreateBook types.RequestCreateBook
		if err := json.NewDecoder(r.Body).Decode(&requestCreateBook); err != nil {
			slog.Error("failed to read request body", slog.Any("error", err))
			writeJSON(w, http.StatusBadRequest, httpError{
				Message: "invalid request body",
			})
			return
		}

		createdBook, err := libraryService.CreateBook(ctx, requestCreateBook.ToBook())
		if err != nil {
			slog.Error("failed to create book", slog.Any("error", err))
			writeJSON(w, http.StatusInternalServerError, httpError{
				Message: "failed to create book",
			})
			return
		}

		writeJSON(w, http.StatusCreated, map[string]any{
			"data": createdBook,
		})
	}
}
