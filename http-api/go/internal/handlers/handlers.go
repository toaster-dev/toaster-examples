package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"toasterexample/internal/services"
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

