package handlers

import (
	"capturecraft-api/internal/middleware"
	"capturecraft-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateShareLink issues a public token for a project.
func (h *Handler) CreateShareLink(c *fiber.Ctx) error {
	userID := middleware.UserID(c)
	projectID := c.Params("projectId")
	project, found := h.Store.FindProjectByID(projectID)
	if !found {
		return fiber.NewError(fiber.StatusNotFound, "project not found")
	}
	if project.OwnerID != userID {
		return fiber.NewError(fiber.StatusForbidden, "not your project")
	}
	link, err := h.Store.CreateShareLink(models.ShareLink{
		Token:     uuid.NewString(),
		ProjectID: projectID,
		OwnerID:   userID,
		CreatedAt: h.now(),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(link)
}
func (h *Handler) GetSharedProject(c *fiber.Ctx) error {
	token := c.Params("token")
	link, found := h.Store.FindShareLinkByToken(token)
	if !found {
		return fiber.NewError(fiber.StatusNotFound, "share link not found")
	}
	project, ok := h.Store.FindProjectByID(link.ProjectID)
	if !ok {
		return fiber.NewError(fiber.StatusNotFound, "project not found")
	}
	screenshots, err := h.Store.ListScreenshotsByProject(project.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(fiber.Map{
		"project":     project,
		"screenshots": screenshots,
	})
}
