package usecase

import (
	"github.com/hydralinkapp/hydralink/api/internal/repository"
	"github.com/hydralinkapp/hydralink/api/pkg/database"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Session represents all methods that Session usecase must implement
type Session interface {
	// Login allows user to login with email and password
	Login(email, password string) error
}

// NewUser creates new User usecase
func NewSession(db *database.DB) Session {
	return sessionImpl{repo: repository.NewUser(db)}
}

type sessionImpl struct {
	repo repository.User
}

// Login allows user to login with email and password
func (s sessionImpl) Login(email, password string) error {
	return AuthenticateUser(s.repo, email, password)
}

// AuthenticateUser checks if given email and password can be authenticated
func AuthenticateUser(repo repository.User, email, password string) error {
	// Find user by email
	found, err := repo.Find(repository.UserSearch{
		Email: email,
	}, false)
	if err != nil {
		return err
	}

	if found.ID == "" {
		return errors.WithStack(ErrEmailInvalid)
	}

	// Attempt to match hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(found.HashedPassword), []byte(password)); err != nil {
		return errors.WithStack(ErrPasswordInvalid)
	}
	return nil
}
