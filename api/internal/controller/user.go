package controller

import (
	"errors"
	"net/http"
	"path/filepath"

	"github.com/hydralinkapp/hydralink/api/internal/entity"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/hydralinkapp/hydralink/api/internal/usecase"
	"github.com/hydralinkapp/hydralink/api/pkg/database"
	"github.com/hydralinkapp/hydralink/api/pkg/monitor"
)

const users = "users"

// changePwdRequest represents request payload
type changePwdRequest struct {
	Email       string `json:"email" binding:"required"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password"  binding:"required"`
}

// signUpRequest represents request payload
type signUpRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password"  binding:"required"`
}

// userResponse represents response payload
type userResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// Session sets up sesion url with controller
func User(router *gin.Engine, logger *monitor.Logger, db *database.DB) {
	u := newUser(logger, db)
	r := router.Group(filepath.Join(root, users))
	r.POST("/signup", u.signUp)
	r.POST("/password/change", u.changePassword)
}

// newHealth creates new health handler
func newUser(logger *monitor.Logger, db *database.DB) *user {
	return &user{
		logger:  logger,
		usecase: usecase.NewUser(db),
	}
}

type user struct {
	logger  *monitor.Logger
	usecase usecase.User
}

// signUp allows user to sign up new account
func (u user) signUp(c *gin.Context) {
	var req signUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing mandatory fields"})
		return
	}

	// TODO Validate req - check email format

	// Generate hashed password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	user, err := u.usecase.SignUp(entity.User{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		HashedPassword: string(hashedPwd),
	})
	if err != nil {
		if errors.Is(err, usecase.ErrUserExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, userResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.FirstName,
		Email:     user.Email,
	})
}

// changePassword allows user to change password
func (u user) changePassword(c *gin.Context) {
	var req changePwdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing mandatory fields"})
		return
	}

	// TODO - add password policy check here

	if err := u.usecase.ChangePassword(req.Email, req.OldPassword, req.NewPassword); err != nil {
		if errors.Is(err, usecase.ErrPasswordInvalid) || errors.Is(err, usecase.ErrEmailInvalid) {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
