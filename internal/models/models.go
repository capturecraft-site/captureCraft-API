package models

import "time"

// User represents an account that owns projects and screenshots.
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
}

// Project groups screenshots under a name/description.
type Project struct {
	ID          string    `json:"id"`
	OwnerID     string    `json:"ownerId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Screenshot is a single uploaded or linked image belonging to a project.
type Screenshot struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"projectId"`
	OwnerID   string    `json:"ownerId"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Status    string    `json:"status"` // e.g. draft, published
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Comment adds discussion on a screenshot.
type Comment struct {
	ID           string    `json:"id"`
	ScreenshotID string    `json:"screenshotId"`
	AuthorID     string    `json:"authorId"`
	Body         string    `json:"body"`
	CreatedAt    time.Time `json:"createdAt"`
}

// ShareLink is a simple public token that surfaces a project.
type ShareLink struct {
	Token     string    `json:"token"`
	ProjectID string    `json:"projectId"`
	OwnerID   string    `json:"ownerId"`
	CreatedAt time.Time `json:"createdAt"`
}
