package domain

import "errors"

// Sentinel Errors
var (
	// ErrNotFound reutilizable
	ErrNotFound = errors.New("recurso no encontrado")
	// ErrAlreadyExists reutilizable
	ErrAlreadyExists = errors.New("el resurso no existe")
	// ErrUnauthorized reutilizable
	ErrUnauthorized = errors.New("no autorizado")
	// ErrForbidden reutilizable
	ErrForbidden = errors.New("acceso prohibido")
	// ErrInvalidInput reutilizable
	ErrInvalidInput = errors.New("datos de entrada inválidos")
	// ErrInternal reutilizable
	ErrInternal = errors.New("error interno")
)
