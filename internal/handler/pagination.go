// Package handler contiene los manejadores HTTP de la aplicación.
package handler

import (
	"net/http"
	"strconv"
)

// Crear un contrato para nuestros paginadores

// PaginationParams representa parámetros de paginación
type PaginationParams struct {
	Cursor string
	Limit  int
}

const (
	defaultLimit = 20
	maxLimit     = 100
)

// ParsePagination lee cursor y limit de la query string
func ParsePagination(request *http.Request) PaginationParams {
	cursor := request.URL.Query().Get("cursor")
	limit := defaultLimit

	if limitQuery := request.URL.Query().Get("limit"); limitQuery != "" {
		if parsed, err := strconv.Atoi(limitQuery); err == nil && parsed > 0 {
			limit = parsed
		} //  strconv.Atoi(limitQuery) convierte un string a entero
	}

	if limit > maxLimit {
		limit = maxLimit
	}

	return PaginationParams{Cursor: cursor, Limit: limit}
}

// PaginatedResponse es el formato estandar para cualquier listado paginado
type PaginatedResponse[T any] struct {
	Data       []T    `json:"data"`
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
}
