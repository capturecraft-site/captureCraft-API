package storage

import "capturecraft-api/internal/models"

// UserStore manages user persistence.
type UserStore interface {
	CreateUser(user models.User) (models.User, error)
	FindUserByEmail(email string) (models.User, bool)
	FindUserByID(id string) (models.User, bool)
}

// ProjectStore manages projects.
type ProjectStore interface {
	CreateProject(project models.Project) (models.Project, error)
	ListProjectsByOwner(ownerID string) ([]models.Project, error)
	FindProjectByID(id string) (models.Project, bool)
	UpdateProject(project models.Project) (models.Project, error)
	DeleteProject(id string) error
}

// ScreenshotStore manages screenshots.
type ScreenshotStore interface {
	CreateScreenshot(s models.Screenshot) (models.Screenshot, error)
	ListScreenshotsByProject(projectID string) ([]models.Screenshot, error)
	FindScreenshotByID(id string) (models.Screenshot, bool)
	UpdateScreenshot(s models.Screenshot) (models.Screenshot, error)
	DeleteScreenshot(id string) error
}

// CommentStore manages comments.
type CommentStore interface {
	CreateComment(c models.Comment) (models.Comment, error)
	ListCommentsByScreenshot(screenshotID string) ([]models.Comment, error)
}

// ShareLinkStore manages share tokens.
type ShareLinkStore interface {
	CreateShareLink(link models.ShareLink) (models.ShareLink, error)
	FindShareLinkByToken(token string) (models.ShareLink, bool)
}

// Store aggregates all store behaviors.
type Store interface {
	UserStore
	ProjectStore
	ScreenshotStore
	CommentStore
	ShareLinkStore
}
