package handlers

import (
	"capturecraft-api/internal/models"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func sanitizeUser(u models.User) models.User {
	u.PasswordHash = ""
	return u
}

// Register creates a new user and returns a JWT.
func (h *Handler) Register(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Name = strings.TrimSpace(req.Name)
	if req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "email and password required")
	}
	if _, exists := h.Store.FindUserByEmail(req.Email); exists {
		return fiber.NewError(fiber.StatusConflict, "email already registered")
	}
	hash, err := h.Auth.HashPassword(req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to hash password")
	}
	user := models.User{
		ID:           uuid.NewString(),
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: hash,
		CreatedAt:    h.now(),
	}
	user, err = h.Store.CreateUser(user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	token, err := h.Auth.GenerateToken(user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to generate token")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"token": token,
		"user":  sanitizeUser(user),
	})
}

// Login authenticates a user and returns a JWT.
func (h *Handler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
	}
	user, found := h.Store.FindUserByEmail(strings.TrimSpace(strings.ToLower(req.Email)))
	if !found {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}
	if err := h.Auth.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}
	token, err := h.Auth.GenerateToken(user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to generate token")
	}
	return c.JSON(fiber.Map{
		"token": token,
		"user":  sanitizeUser(user),
	})
}
