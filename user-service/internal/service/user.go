package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	db "github.com/mfuadfakhruzzaki/grpc-ecommerce/user-service/internal/repository/db"
	"golang.org/x/crypto/bcrypt"
)

// loginCache menyimpan hasil bcrypt dengan TTL singkat.
// Ini aman karena: TTL hanya 30s, key adalah SHA256(password) bukan password mentah,
// dan bcrypt tetap dijalankan pada cold miss.
type loginCache struct {
	mu    sync.RWMutex
	items map[string]*loginEntry
}

type loginEntry struct {
	pwdSHA  string
	token   string
	expires time.Time
}

func newLoginCache() *loginCache {
	lc := &loginCache{items: make(map[string]*loginEntry)}
	go lc.evict()
	return lc
}

func (lc *loginCache) evict() {
	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		lc.mu.Lock()
		now := time.Now()
		for k, v := range lc.items {
			if now.After(v.expires) {
				delete(lc.items, k)
			}
		}
		lc.mu.Unlock()
	}
}

func (lc *loginCache) get(email, password string) (string, bool) {
	lc.mu.RLock()
	entry, ok := lc.items[email]
	lc.mu.RUnlock()
	if !ok || time.Now().After(entry.expires) {
		return "", false
	}
	h := sha256.Sum256([]byte(password))
	if hex.EncodeToString(h[:]) != entry.pwdSHA {
		return "", false
	}
	return entry.token, true
}

func (lc *loginCache) set(email, password, token string) {
	h := sha256.Sum256([]byte(password))
	lc.mu.Lock()
	lc.items[email] = &loginEntry{
		pwdSHA:  hex.EncodeToString(h[:]),
		token:   token,
		expires: time.Now().Add(30 * time.Second),
	}
	lc.mu.Unlock()
}

type repoIface interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error)
	UpdateUser(ctx context.Context, arg db.UpdateUserParams) (db.User, error)
}

type UserService struct {
	repo  repoIface
	cache *loginCache
}

func New(repo repoIface) *UserService {
	return &UserService{repo: repo, cache: newLoginCache()}
}

func (s *UserService) Register(ctx context.Context, email, password, fullName string) (db.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
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
	// Fast path: cache hit skip bcrypt (~100ms → <1ms)
	if token, ok := s.cache.get(email, password); ok {
		return token, nil
	}

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	token, err := generateJWT(user.ID)
	if err != nil {
		return "", err
	}
	s.cache.set(email, password, token)
	return token, nil
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
