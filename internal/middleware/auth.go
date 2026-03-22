package middleware

import (
	"capturecraft-api/internal/auth"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// RequireAuth validates JWT bearer tokens and injects user ID into context.
func RequireAuth(authSvc *auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" || !strings.HasPrefix(strings.ToLower(header), "bearer ") {
			return fiber.NewError(fiber.StatusUnauthorized, "missing bearer token")
		}
		token := strings.TrimSpace(header[7:])
		userID, err := authSvc.ParseToken(token)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}
		c.Locals("userID", userID)
		return c.Next()
	}
}

// UserID extracts the authenticated user ID from context.
func UserID(c *fiber.Ctx) string {
	if v := c.Locals("userID"); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
