package usecase

import (
	"github.com/hydralinkapp/hydralink/api/internal/repository"
	"github.com/hydralinkapp/hydralink/api/pkg/database"
)

// Health wraps repository
type Health struct {
	repo repository.Health
}

// NewHealth creates new health usecase
func NewHealth(db *database.DB) Health {
	return Health{repo: repository.NewHealth(db)}
}

// readines checks for ok database connection
func (h Health) Readiness() error {
	return h.repo.Readiness()
}
