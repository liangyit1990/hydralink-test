package repository

import (
	"gorm.io/gorm"

	"github.com/hydralinkapp/hydralink/api/internal/entity"
	"github.com/hydralinkapp/hydralink/api/pkg/database"
	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

// NewUser creates new user impl and returns as User interface
func NewUser(db *database.DB) User {
	return &userImpl{db: db}
}

// User represents methods that User repository must implement
type User interface {
	// DB returns the db
	DB() *database.DB
	// Create inserts new record in User table
	Create(u entity.User) (entity.User, error)
	// Find retrieves a user based on search criteria
	Find(s UserSearch, lock bool) (entity.User, error)
	// Update retrieves user for locking and then update
	Update(u entity.User) error
}

// UserSearch sets fields to search user
type UserSearch struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
}

// UserUpdate sets fields to update user
type UserUpdate struct {
	HashedPassword string
}

type userImpl struct {
	db *database.DB
}

// DB returns the db
func (u userImpl) DB() *database.DB {
	return u.db
}

// Create inserts new record in User table
func (u userImpl) Create(user entity.User) (entity.User, error) {
	if err := u.db.GormDB.Create(&user).Error; err != nil {
		return entity.User{}, errors.WithStack(err)
	}
	return user, nil
}

// Find retrieves an user based on search criteria
func (u userImpl) Find(filter UserSearch, lock bool) (entity.User, error) {
	tx := u.db.GormDB
	// lock basically locks the row preventing others from performing update
	if lock {
		tx = tx.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	if filter.ID != "" {
		tx = tx.Where("id = ?", filter.ID)
	}

	if filter.Email != "" {
		tx = tx.Where("email = ?", filter.Email)
	}

	if filter.FirstName != "" {
		tx = tx.Where("first_name = ?", filter.FirstName)
	}

	if filter.LastName != "" {
		tx = tx.Where("last_name = ?", filter.LastName)
	}

	var result entity.User
	if err := tx.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, nil
		}
		return entity.User{}, errors.WithStack(err)
	}

	return result, nil
}

// Update retrieves user for locking and then update
func (u userImpl) Update(user entity.User) error {
	if user.ID == "" {
		return errors.WithStack(errors.New("user to be update is empty"))
	}
	return errors.WithStack(u.db.GormDB.Save(&user).Error)
}
