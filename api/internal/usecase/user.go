package usecase

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sony/sonyflake"
	"golang.org/x/crypto/bcrypt"

	"github.com/hydralinkapp/hydralink/api/internal/entity"
	"github.com/hydralinkapp/hydralink/api/internal/repository"
	"github.com/hydralinkapp/hydralink/api/pkg/database"
)

// User represents all methods that User usecase must implement
type User interface {
	SignUp(entity.User) (entity.User, error)
	ChangePassword(email string, oldPassword, NewPassword string) error
}

// NewUser creates new User usecase
func NewUser(db *database.DB) User {
	return userImpl{repo: repository.NewUser(db)}
}

type userImpl struct {
	repo repository.User
}

// SignUp signs up a new user
func (u userImpl) SignUp(in entity.User) (entity.User, error) {
	// Check if user by first and last name exists
	found, err := u.repo.Find(repository.UserSearch{
		FirstName: in.FirstName,
		LastName:  in.LastName,
	}, false)
	if err != nil {
		return entity.User{}, errors.WithStack(err)
	}

	if found.ID != "" {
		return entity.User{}, errors.WithStack(ErrUserExists)
	}

	// Generate unique id
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		return entity.User{}, errors.WithStack(err)
	}
	in.ID = fmt.Sprintf("%x", id)

	// No existing user, proceed to create user
	return u.repo.Create(in)
}

// ChangePassword changes the user password
func (u userImpl) ChangePassword(email string, oldPassword, NewPassword string) error {
	// Check if user is account owner
	if err := AuthenticateUser(u.repo, email, oldPassword); err != nil {
		return err
	}

	// Retrieve user entry with lock
	locked, err := u.repo.Find(repository.UserSearch{
		Email: email,
	}, true)
	if err != nil {
		return err
	}

	hashedNewPwd, err := bcrypt.GenerateFromPassword([]byte(NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.WithStack(err)
	}

	locked.HashedPassword = string(hashedNewPwd)

	return u.repo.Update(locked)
}
