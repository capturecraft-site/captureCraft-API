package handlers

import (
	"capturecraft-api/internal/middleware"
	"capturecraft-api/internal/models"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateScreenshot creates a screenshot under a project.
func (h *Handler) CreateScreenshot(c *fiber.Ctx) error {
	userID := middleware.UserID(c)
	projectID := c.Params("projectId")
	project, found := h.Store.FindProjectByID(projectID)
	if !found {
		return fiber.NewError(fiber.StatusNotFound, "project not found")
	}
	if project.OwnerID != userID {
		return fiber.NewError(fiber.StatusForbidden, "not your project")
	}
	var req struct {
		Title  string `json:"title"`
		URL    string `json:"url"`
		Status string `json:"status"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
	}
	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = "draft"
	}
	screenshot := models.Screenshot{
		ID:        uuid.NewString(),
		ProjectID: projectID,
		OwnerID:   userID,
		Title:     strings.TrimSpace(req.Title),
		URL:       strings.TrimSpace(req.URL),
		Status:    status,
		CreatedAt: h.now(),
		UpdatedAt: h.now(),
	}
	created, err := h.Store.CreateScreenshot(screenshot)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(created)
}

// ListScreenshots lists screenshots for a project.
func (h *Handler) ListScreenshots(c *fiber.Ctx) error {
	userID := middleware.UserID(c)
	projectID := c.Params("projectId")
	project, found := h.Store.FindProjectByID(projectID)
	if !found {
		return fiber.NewError(fiber.StatusNotFound, "project not found")
	}
	if project.OwnerID != userID {
		return fiber.NewError(fiber.StatusForbidden, "not your project")
	}
	items, err := h.Store.ListScreenshotsByProject(projectID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(items)
}

// GetScreenshot fetches a screenshot by ID.
func (h *Handler) GetScreenshot(c *fiber.Ctx) error {
	userID := middleware.UserID(c)
	id := c.Params("id")
	sc, found := h.Store.FindScreenshotByID(id)
	if !found {
		return fiber.NewError(fiber.StatusNotFound, "screenshot not found")
	}
	if sc.OwnerID != userID {
		return fiber.NewError(fiber.StatusForbidden, "not your screenshot")
	}
	return c.JSON(sc)
}

// UpdateScreenshot updates title, url, or status.
func (h *Handler) UpdateScreenshot(c *fiber.Ctx) error {
	userID := middleware.UserID(c)
	id := c.Params("id")
	sc, found := h.Store.FindScreenshotByID(id)
	if !found {
		return fiber.NewError(fiber.StatusNotFound, "screenshot not found")
	}
	if sc.OwnerID != userID {
		return fiber.NewError(fiber.StatusForbidden, "not your screenshot")
	}
	var req struct {
		Title  *string `json:"title"`
		URL    *string `json:"url"`
		Status *string `json:"status"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
	}
	if req.Title != nil {
		sc.Title = strings.TrimSpace(*req.Title)
	}
	if req.URL != nil {
		sc.URL = strings.TrimSpace(*req.URL)
	}
	if req.Status != nil {
		status := strings.TrimSpace(*req.Status)
		if status != "" {
			sc.Status = status
		}
	}
	sc.UpdatedAt = h.now()
	updated, err := h.Store.UpdateScreenshot(sc)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(updated)
}

// DeleteScreenshot removes a screenshot.
func (h *Handler) DeleteScreenshot(c *fiber.Ctx) error {
	userID := middleware.UserID(c)
	id := c.Params("id")
	sc, found := h.Store.FindScreenshotByID(id)
	if !found {
		return fiber.NewError(fiber.StatusNotFound, "screenshot not found")
	}
	if sc.OwnerID != userID {
		return fiber.NewError(fiber.StatusForbidden, "not your screenshot")
	}
	if err := h.Store.DeleteScreenshot(id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
