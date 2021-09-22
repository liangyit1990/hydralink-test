package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hydralinkapp/hydralink/api/internal/usecase"
	"github.com/hydralinkapp/hydralink/api/pkg/database"
	"github.com/hydralinkapp/hydralink/api/pkg/monitor"
)

const root = "api/v1"

// Health sets up health url with controller
func Health(router *gin.Engine, logger *monitor.Logger, db *database.DB) {
	h := newHealth(logger, db)
	r := router.Group(root)
	r.GET("/liveness", h.liveness)
	r.GET("/readiness", h.readiness)
}

// newHealth creates new health handler
func newHealth(logger *monitor.Logger, db *database.DB) *health {
	return &health{
		logger:  logger,
		usecase: usecase.NewHealth(db),
	}
}

type health struct {
	logger  *monitor.Logger
	usecase usecase.Health
}

// liveness checks for ok connection
func (h health) liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Live!",
	})
}

// readiness checks for ok database connection
func (h health) readiness(c *gin.Context) {
	if err := h.usecase.Readiness(); err != nil {
		h.logger.Errorf("%+v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Ready!",
	})
}
