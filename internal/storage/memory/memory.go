package memory

import (
	"capturecraft-api/internal/models"
	"capturecraft-api/internal/storage"
	"errors"
	"sync"
)

// Store is an in-memory implementation of storage.Store.
type Store struct {
	mu          sync.RWMutex
	users       map[string]models.User
	projects    map[string]models.Project
	screenshots map[string]models.Screenshot
	comments    map[string]models.Comment
	shareLinks  map[string]models.ShareLink
}

// New creates a new in-memory store.
func New() *Store {
	return &Store{
		users:       make(map[string]models.User),
		projects:    make(map[string]models.Project),
		screenshots: make(map[string]models.Screenshot),
		comments:    make(map[string]models.Comment),
		shareLinks:  make(map[string]models.ShareLink),
	}
}

// Ensure Store satisfies storage.Store.
var _ storage.Store = (*Store)(nil)

// User operations
func (s *Store) CreateUser(user models.User) (models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.users[user.ID]; exists {
		return models.User{}, errors.New("user already exists")
	}
	s.users[user.ID] = user
	return user, nil
}

func (s *Store) FindUserByEmail(email string) (models.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.users {
		if u.Email == email {
			return u, true
		}
	}
	return models.User{}, false
}

func (s *Store) FindUserByID(id string) (models.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	return u, ok
}

// Project operations
func (s *Store) CreateProject(project models.Project) (models.Project, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.projects[project.ID]; exists {
		return models.Project{}, errors.New("project already exists")
	}
	s.projects[project.ID] = project
	return project, nil
}

func (s *Store) ListProjectsByOwner(ownerID string) ([]models.Project, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var projects []models.Project
	for _, p := range s.projects {
		if p.OwnerID == ownerID {
			projects = append(projects, p)
		}
	}
	return projects, nil
}

func (s *Store) FindProjectByID(id string) (models.Project, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.projects[id]
	return p, ok
}

func (s *Store) UpdateProject(project models.Project) (models.Project, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.projects[project.ID]; !exists {
		return models.Project{}, errors.New("project not found")
	}
	s.projects[project.ID] = project
	return project, nil
}

func (s *Store) DeleteProject(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.projects, id)
	return nil
}

// Screenshot operations
func (s *Store) CreateScreenshot(sc models.Screenshot) (models.Screenshot, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.screenshots[sc.ID]; exists {
		return models.Screenshot{}, errors.New("screenshot already exists")
	}
	s.screenshots[sc.ID] = sc
	return sc, nil
}

func (s *Store) ListScreenshotsByProject(projectID string) ([]models.Screenshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []models.Screenshot
	for _, sc := range s.screenshots {
		if sc.ProjectID == projectID {
			result = append(result, sc)
		}
	}
	return result, nil
}

func (s *Store) FindScreenshotByID(id string) (models.Screenshot, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sc, ok := s.screenshots[id]
	return sc, ok
}

func (s *Store) UpdateScreenshot(sc models.Screenshot) (models.Screenshot, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.screenshots[sc.ID]; !exists {
		return models.Screenshot{}, errors.New("screenshot not found")
	}
	s.screenshots[sc.ID] = sc
	return sc, nil
}

func (s *Store) DeleteScreenshot(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.screenshots, id)
	return nil
}

// Comment operations
func (s *Store) CreateComment(c models.Comment) (models.Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.comments[c.ID]; exists {
		return models.Comment{}, errors.New("comment already exists")
	}
	s.comments[c.ID] = c
	return c, nil
}

func (s *Store) ListCommentsByScreenshot(screenshotID string) ([]models.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []models.Comment
	for _, c := range s.comments {
		if c.ScreenshotID == screenshotID {
			result = append(result, c)
		}
	}
	return result, nil
}

// ShareLink operations
func (s *Store) CreateShareLink(link models.ShareLink) (models.ShareLink, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.shareLinks[link.Token] = link
	return link, nil
}

func (s *Store) FindShareLinkByToken(token string) (models.ShareLink, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	l, ok := s.shareLinks[token]
	return l, ok
}
