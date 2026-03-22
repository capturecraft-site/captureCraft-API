package auth

import (
	"capturecraft-api/internal/models"
	"capturecraft-api/internal/storage"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Service handles auth-related helpers.
type Service struct {
	secret   []byte
	tokenTTL time.Duration
	store    storage.UserStore
}

// NewService builds a Service.
func NewService(secret string, tokenTTL time.Duration, store storage.UserStore) *Service {
	return &Service{secret: []byte(secret), tokenTTL: tokenTTL, store: store}
}

// HashPassword hashes a plaintext password.
func (s *Service) HashPassword(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compares plaintext to a hash.
func (s *Service) CheckPassword(hash, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}

// GenerateToken creates a signed JWT for a user.
func (s *Service) GenerateToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(s.tokenTTL).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(s.secret)
}

// ParseToken validates a token and returns the subject (user ID).
func (s *Service) ParseToken(token string) (string, error) {
	tok, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := tok.Claims.(jwt.MapClaims); ok && tok.Valid {
		sub, _ := claims["sub"].(string)
		if sub == "" {
			return "", errors.New("missing subject")
		}
		return sub, nil
	}
	return "", errors.New("invalid token")
}
