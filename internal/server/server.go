package server

import (
	"capturecraft-api/internal/auth"
	"capturecraft-api/internal/handlers"
	"capturecraft-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// New builds the Fiber app with all routes and middleware.
func New(handler *handlers.Handler, authSvc *auth.Service) *fiber.App {
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	api := app.Group("/api/v1")
	api.Post("/auth/register", handler.Register)
	api.Post("/auth/login", handler.Login)

	protected := api.Group("")
	protected.Use(middleware.RequireAuth(authSvc))

	// Projects
	protected.Post("/projects", handler.CreateProject)
	protected.Get("/projects", handler.ListProjects)
	protected.Get("/projects/:id", handler.GetProject)
	protected.Patch("/projects/:id", handler.UpdateProject)
	protected.Delete("/projects/:id", handler.DeleteProject)

	// Screenshots
	protected.Post("/projects/:projectId/screenshots", handler.CreateScreenshot)
	protected.Get("/projects/:projectId/screenshots", handler.ListScreenshots)
	protected.Get("/screenshots/:id", handler.GetScreenshot)
	protected.Patch("/screenshots/:id", handler.UpdateScreenshot)
	protected.Delete("/screenshots/:id", handler.DeleteScreenshot)

	// Comments
	protected.Get("/screenshots/:screenshotId/comments", handler.ListComments)
	protected.Post("/screenshots/:screenshotId/comments", handler.CreateComment)

	// Share links
	protected.Post("/projects/:projectId/share", handler.CreateShareLink)

	// Uploads
	protected.Post("/uploads/presign", handler.PresignUpload)

	// Public share endpoint
	public := app.Group("/public")
	public.Get("/share/:token", handler.GetSharedProject)

	return app
}
