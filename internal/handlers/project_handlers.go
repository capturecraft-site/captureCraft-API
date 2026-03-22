package handlers

import (
	"capturecraft-api/internal/middleware"
	"capturecraft-api/internal/models"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateProject makes a new project for the authenticated user.
func (h *Handler) CreateProject(c *fiber.Ctx) error {
	userID := middleware.UserID(c)
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "name is required")
	}
	project := models.Project{
		ID:          uuid.NewString(),
		OwnerID:     userID,
		Name:        req.Name,
		Description: strings.TrimSpace(req.Description),
		CreatedAt:   h.now(),
		UpdatedAt:   h.now(),
	}
	created, err := h.Store.CreateProject(project)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(created)
}

// ListProjects lists current user's projects.
func (h *Handler) ListProjects(c *fiber.Ctx) error {
	userID := middleware.UserID(c)
	projects, err := h.Store.ListProjectsByOwner(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(projects)
}

// GetProject fetches a single project by ID if owned by the user.
func (h *Handler) GetProject(c *fiber.Ctx) error {
	userID := middleware.UserID(c)
	id := c.Params("id")
	project, found := h.Store.FindProjectByID(id)
	if !found {
		return fiber.NewError(fiber.StatusNotFound, "project not found")
	}
	if project.OwnerID != userID {
		return fiber.NewError(fiber.StatusForbidden, "not your project")
	}
	return c.JSON(project)
}

// UpdateProject patches project fields.
func (h *Handler) UpdateProject(c *fiber.Ctx) error {
	userID := middleware.UserID(c)
	id := c.Params("id")
	project, found := h.Store.FindProjectByID(id)
	if !found {
		return fiber.NewError(fiber.StatusNotFound, "project not found")
	}
	if project.OwnerID != userID {
		return fiber.NewError(fiber.StatusForbidden, "not your project")
	}
	var req struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
	}
	if req.Name != nil {
		project.Name = strings.TrimSpace(*req.Name)
	}
	if req.Description != nil {
		project.Description = strings.TrimSpace(*req.Description)
	}
	project.UpdatedAt = h.now()
	updated, err := h.Store.UpdateProject(project)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(updated)
}

// DeleteProject removes a project.
func (h *Handler) DeleteProject(c *fiber.Ctx) error {
	userID := middleware.UserID(c)
	id := c.Params("id")
	project, found := h.Store.FindProjectByID(id)
	if !found {
		return fiber.NewError(fiber.StatusNotFound, "project not found")
	}
	if project.OwnerID != userID {
		return fiber.NewError(fiber.StatusForbidden, "not your project")
	}
	if err := h.Store.DeleteProject(id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
