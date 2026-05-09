package admin

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Service interface {
	CreateAdmin(ctx context.Context, req CreateAdminRequest) (AdminResponse, error)
	LoginAdmin(ctx context.Context, req AdminLoginRequest) (AdminLoginResponse, error)
	GetAdmin(ctx context.Context, id uuid.UUID) (AdminResponse, error)
	UpdateAdmin(ctx context.Context, id uuid.UUID, req UpdateAdminRequest) error
	DeleteAdmin(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repository
	jwt  *JWTManager
}

func NewService(repo Repository, jwtManager *JWTManager) Service {
	return &service{
		repo: repo,
		jwt:  jwtManager,
	}
}

func (s *service) CreateAdmin(ctx context.Context, req CreateAdminRequest) (AdminResponse, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return AdminResponse{}, err
	}

	admin, err := s.repo.CreateAdmin(ctx, req.Name, req.Email, string(passwordHash))
	if err != nil {
		return AdminResponse{}, err
	}

	return AdminResponse{
		ID:    admin.ID.String(),
		Name:  admin.Name,
		Email: admin.Email,
	}, nil
}

func (s *service) LoginAdmin(ctx context.Context, req AdminLoginRequest) (AdminLoginResponse, error) {
	admin, err := s.repo.GetAdminByEmail(ctx, req.Email)
	if err != nil {
		return AdminLoginResponse{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(req.Password)); err != nil {
		return AdminLoginResponse{}, ErrInvalidCredentials
	}

	token, expiresAt, err := s.jwt.GenerateToken(admin.ID.String(), admin.Email, admin.IsSuperadmin)
	if err != nil {
		return AdminLoginResponse{}, err
	}

	return AdminLoginResponse{
		Token:     token,
		ExpiresAt: expiresAt.Unix(),
	}, nil
}

func (s *service) GetAdmin(ctx context.Context, id uuid.UUID) (AdminResponse, error) {
	admin, err := s.repo.GetAdminByID(ctx, id)
	if err != nil {
		return AdminResponse{}, err
	}

	return AdminResponse{
		ID:    admin.ID.String(),
		Name:  admin.Name,
		Email: admin.Email,
	}, nil
}

func (s *service) UpdateAdmin(ctx context.Context, id uuid.UUID, req UpdateAdminRequest) error {
	return s.repo.UpdateAdmin(ctx, id, req)
}

func (s *service) DeleteAdmin(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteAdmin(ctx, id)
}
