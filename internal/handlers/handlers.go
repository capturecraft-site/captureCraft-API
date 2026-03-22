package handlers

import (
	"capturecraft-api/internal/auth"
	"capturecraft-api/internal/storage"
	"time"
)

// Handler bundles dependencies for HTTP handlers.
type Handler struct {
	Store storage.Store
	Auth  *auth.Service
}

// New constructs a Handler.
func New(store storage.Store, authSvc *auth.Service) *Handler {
	return &Handler{Store: store, Auth: authSvc}
}

// now() is kept here to simplify testing/mocking later.
func (h *Handler) now() time.Time { return time.Now().UTC() }
