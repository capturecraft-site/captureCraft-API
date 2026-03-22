package handlers

import (
	"capturecraft-api/internal/middleware"
	"capturecraft-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ListComments returns comments for a screenshot if caller owns it.
func (h *Handler) ListComments(c *fiber.Ctx) error {
	userID := middleware.UserID(c)
	screenshotID := c.Params("screenshotId")
	sc, found := h.Store.FindScreenshotByID(screenshotID)
	if !found {
		return fiber.NewError(fiber.StatusNotFound, "screenshot not found")
	}
	if sc.OwnerID != userID {
		return fiber.NewError(fiber.StatusForbidden, "not your screenshot")
	}
	items, err := h.Store.ListCommentsByScreenshot(screenshotID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(items)
}

// CreateComment adds a comment to a screenshot.
func (h *Handler) CreateComment(c *fiber.Ctx) error {
	userID := middleware.UserID(c)
	screenshotID := c.Params("screenshotId")
	sc, found := h.Store.FindScreenshotByID(screenshotID)
	if !found {
		return fiber.NewError(fiber.StatusNotFound, "screenshot not found")
	}
	if sc.OwnerID != userID {
		return fiber.NewError(fiber.StatusForbidden, "not your screenshot")
	}
	var req struct {
		Body string `json:"body"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
	}
	if req.Body == "" {
		return fiber.NewError(fiber.StatusBadRequest, "body is required")
	}
	comment := models.Comment{
		ID:           uuid.NewString(),
		ScreenshotID: screenshotID,
		AuthorID:     userID,
		Body:         req.Body,
		CreatedAt:    h.now(),
	}
	created, err := h.Store.CreateComment(comment)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(created)
}
