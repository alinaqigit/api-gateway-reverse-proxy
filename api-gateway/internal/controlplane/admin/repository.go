package admin

import (
	"context"

	"api-gateway/internal/db"

	"github.com/google/uuid"
)

type Repository interface {
	CreateAdmin(ctx context.Context, name string, email string, passwordHash string) (db.Admin, error)
	GetAdminByEmail(ctx context.Context, email string) (db.Admin, error)
	GetAdminByID(ctx context.Context, id uuid.UUID) (db.GetAdminByIDRow, error)
	UpdateAdmin(ctx context.Context, id uuid.UUID, req UpdateAdminRequest) error
	DeleteAdmin(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) Repository {
	return &repository{queries: queries}
}

func (r *repository) CreateAdmin(ctx context.Context, name string, email string, passwordHash string) (db.Admin, error) {
	params := db.CreateAdminParams{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	}
	return r.queries.CreateAdmin(ctx, params)
}

func (r *repository) GetAdminByEmail(ctx context.Context, email string) (db.Admin, error) {
	return r.queries.GetAdminByEmail(ctx, email)
}

func (r *repository) GetAdminByID(ctx context.Context, id uuid.UUID) (db.GetAdminByIDRow, error) {
	return r.queries.GetAdminByID(ctx, id)
}

func (r *repository) UpdateAdmin(ctx context.Context, id uuid.UUID, req UpdateAdminRequest) error {
	params := db.UpdateAdminParams{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	}
	return r.queries.UpdateAdmin(ctx, params)
}

func (r *repository) DeleteAdmin(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteAdmin(ctx, id)
}
