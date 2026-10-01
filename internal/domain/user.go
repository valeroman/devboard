package domain

import "time"

// User representa al usuario del sistema.
type User struct {
	ID           string
	Email        string
	PasswordHash string
	Name         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewUser crea una nueva instancia de User con los datos proporcionados.
func NewUser(email, name, passwordHash string) *User {
	now := time.Now()

	return &User{
		Email:        email,
		Name:         name,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// IsValid verifica si el usuario tiene los campos obligatorios completos.
func (u *User) IsValid() bool {
	return u.Email != "" && u.Name != ""
}
