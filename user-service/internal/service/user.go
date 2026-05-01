package service

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	db "github.com/mfuadfakhruzzaki/grpc-ecommerce/user-service/internal/repository/db"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo interface {
		CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
		GetUserByEmail(ctx context.Context, email string) (db.User, error)
		GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error)
		UpdateUser(ctx context.Context, arg db.UpdateUserParams) (db.User, error)
	}
}

func New(repo interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error)
	UpdateUser(ctx context.Context, arg db.UpdateUserParams) (db.User, error)
}) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, email, password, fullName string) (db.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return db.User{}, err
	}
	return s.repo.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: string(hash),
		FullName:     sql.NullString{String: fullName, Valid: fullName != ""},
	})
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	return generateJWT(user.ID)
}

func (s *UserService) GetProfile(ctx context.Context, userID uuid.UUID) (db.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, fullName, avatarURL string) (db.User, error) {
	return s.repo.UpdateUser(ctx, db.UpdateUserParams{
		ID:        userID,
		FullName:  sql.NullString{String: fullName, Valid: fullName != ""},
		AvatarUrl: sql.NullString{String: avatarURL, Valid: avatarURL != ""},
	})
}

func generateJWT(userID uuid.UUID) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}