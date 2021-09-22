package controller

import (
	"errors"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/hydralinkapp/hydralink/api/internal/usecase"
	"github.com/hydralinkapp/hydralink/api/pkg/database"
	"github.com/hydralinkapp/hydralink/api/pkg/monitor"
)

const sessions = "sessions"

// login represents request payload
type login struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password"  binding:"required"`
}

// Session sets up sesion url with controller
func Session(router *gin.Engine, logger *monitor.Logger, db *database.DB) {
	s := newSession(logger, db)
	r := router.Group(filepath.Join(root, sessions))
	r.POST("/login", s.login)
}

// newHealth creates new health handler
func newSession(logger *monitor.Logger, db *database.DB) *session {
	return &session{
		logger:  logger,
		usecase: usecase.NewSession(db),
	}
}

type session struct {
	logger  *monitor.Logger
	usecase usecase.Session
}

// login allows user to login
func (s session) login(c *gin.Context) {
	var request login
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing email or password"})
		return
	}

	if err := s.usecase.Login(request.Email, request.Password); err != nil {
		if errors.Is(err, usecase.ErrEmailInvalid) || errors.Is(err, usecase.ErrPasswordInvalid) {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
