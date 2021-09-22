package repository

import (
	"github.com/hydralinkapp/hydralink/api/pkg/database"
	"github.com/pkg/errors"
)

type Health struct {
	db *database.DB
}

// NewHealth creates new health controller
func NewHealth(db *database.DB) Health {
	return Health{db: db}
}

// readines checks connection with database
func (h Health) Readiness() error {
	tx := h.db.GormDB
	return errors.WithStack(tx.Exec("SELECT 1").Error)
}
