package usecase

import (
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/hydralinkapp/hydralink/api/internal/entity"
	"github.com/hydralinkapp/hydralink/api/pkg/database"
	"github.com/stretchr/testify/require"
)

func TestUserImpl_SignUp(t *testing.T) {
	type args struct {
		given entity.User
		err   error
	}

	tcs := map[string]args{
		"user already exists": {
			given: mockTestUser,
			err:   ErrPasswordInvalid,
		},
		"success": {
			given: entity.User{
				FirstName:      "newuser",
				LastName:       "newuser",
				HashedPassword: "newuser",
				Email:          "newuser",
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			database.TestTransaction(database.TestDB(), func(dbTx *database.DB) {
				// Seed
				u := NewUser(dbTx)
				dbTx.GormDB.Create(&mockTestUser)
				// When
				_, err := u.SignUp(tc.given)
				// Then
				if tc.err != nil {
					require.Error(t, errors.Unwrap(err), tc.err.Error())
				} else {
					require.NoError(t, err)
				}
			})
		})
	}
}

func TestUserImpl_ChangePassword(t *testing.T) {
	type args struct {
		givenEmail       string
		givenOldPassword string
		givenNewPassword string
		err              error
	}

	tcs := map[string]args{
		"wrong password": {
			givenEmail:       mockTestUser.Email,
			givenOldPassword: "wrongpassword",
			givenNewPassword: "newPassword",
			err:              ErrPasswordInvalid,
		},
		"success": {
			givenEmail:       mockTestUser.Email,
			givenOldPassword: "test",
			givenNewPassword: "newPassword",
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			database.TestTransaction(database.TestDB(), func(dbTx *database.DB) {
				u := NewUser(dbTx)
				dbTx.GormDB.Create(&mockTestUser)
				// When
				err := u.ChangePassword(mockTestUser.Email, tc.givenOldPassword, tc.givenNewPassword)
				// Then
				if tc.err != nil {
					require.Error(t, errors.Unwrap(err), tc.err.Error())
				} else {
					require.NoError(t, err)
					got := entity.User{}
					require.NoError(t, dbTx.GormDB.Find(&got).Error)
					require.NoError(t, bcrypt.CompareHashAndPassword([]byte(got.HashedPassword), []byte(tc.givenNewPassword)))
				}
			})
		})
	}
}
