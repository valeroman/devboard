package memory

import (
	"context"
	"errors"
	"testing"

	"github.com/valeroman/devboard/internal/domain"
)

func TestUserRepository_Create(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	user := &domain.User{
		Name:  "Ana",
		Email: "ana@devboard.dev",
	}

	err := repo.Create(ctx, user)

	if err != nil {
		t.Fatalf("Create() devolvio un error inesperado: %v", err)
	}

	if user.ID == "" {
		t.Fatalf("Create() no asigno un ID al usuario")
	}

	if user.ID != "user_1" {
		t.Fatalf("Create() asigno un ID %q; se esperaba %q", user.ID, "user_1")
	}
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	_, err := repo.GetByID(ctx, "user_999")

	if err == nil {
		t.Fatalf("GetByID()debía devolver un error")
	}

	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("GetByID() devolvió %v; se esperaba ErrNotFound", err)
	}
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	_, err := repo.GetByEmail(ctx, "nonexistent@devboard.dev")

	if err == nil {
		t.Fatalf("GetByEmail() debía devolver un error")
	}

	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("GetByEmail() devolvió %v; se esperaba ErrNotFound", err)
	}
}

// TextUserReposotory_GetByID_Success
func TextUserReposotoryGetByIDSuccess(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	user := &domain.User{
		Name:  "Ana",
		Email: "ana@devboard.dev",
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Create() devolvio un error inesperado: %v", err)
	}

	retrievedUser, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetByID() devolvió un error inesperado: %v", err)
	}

	if retrievedUser.ID != user.ID {
		t.Errorf("GetByID() devolvió un usuario con ID %q; se esperaba %q", retrievedUser.ID, user.ID)
	}
}
