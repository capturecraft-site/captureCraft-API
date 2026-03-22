package handlers

import (
	"capturecraft-api/internal/middleware"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// PresignUpload simulates a presigned URL for uploads (swap with real object storage later).
func (h *Handler) PresignUpload(c *fiber.Ctx) error {
	_ = middleware.UserID(c) // ensures auth middleware ran
	var req struct {
		Filename    string `json:"filename"`
		ContentType string `json:"contentType"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
	}
	req.Filename = strings.TrimSpace(req.Filename)
	if req.Filename == "" {
		return fiber.NewError(fiber.StatusBadRequest, "filename is required")
	}
	u := fmt.Sprintf("https://storage.fake.local/uploads/%s/%s", uuid.NewString(), url.PathEscape(req.Filename))
	return c.JSON(fiber.Map{
		"uploadURL": u,
		"method":    "PUT",
		"headers": fiber.Map{
			"Content-Type": req.ContentType,
		},
		"expiresIn": int((15 * time.Minute).Seconds()),
	})
}
