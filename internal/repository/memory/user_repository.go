// Package memory proporciona una implementación en memoria del repositorio de usuarios.
package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/valeroman/devboard/internal/domain"
)

// UserRepository es un adapter que implementa el contrato UserRepository en memoria.
type UserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

// NewUserRepository (constructor) crea una nueva instancia de UserRepository en memoria.
func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]*domain.User),
	}
}

// Create (metodo) agrega un nuevo usuario al repositorio en memoria.
func (r *UserRepository) Create(_ context.Context, user *domain.User) error {
	r.mu.Lock() // Adquiere un bloqueo de escritura para modificar el mapa de usuarios de manera segura.
	defer r.mu.Unlock()

	user.ID = fmt.Sprintf("user_%d", len(r.users)+1)
	r.users[user.ID] = user
	return nil

}

// GetByID (metodo) obtiene un usuario por su ID desde el repositorio en memoria.
func (r *UserRepository) GetByID(_ context.Context, id string) (*domain.User, error) {
	r.mu.RLock() // Adquiere un bloqueo de lectura para acceder a los usuarios de manera segura.
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("usuario %s: %w", id, domain.ErrNotFound)
	}
	return user, nil
}

// GetByEmail (metodo) obtiene un usuario por su email desde el repositorio en memoria.
func (r *UserRepository) GetByEmail(_ context.Context, email string) (*domain.User, error) {
	r.mu.RLock() // Adquiere un bloqueo de lectura para acceder a los usuarios de manera segura.
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, fmt.Errorf("email %s: %w", email, domain.ErrNotFound)
}
