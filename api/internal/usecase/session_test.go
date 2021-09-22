package usecase

import (
	"testing"

	"github.com/hydralinkapp/hydralink/api/internal/entity"
	"github.com/hydralinkapp/hydralink/api/pkg/database"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

var (
	// global mock user for test sharing
	mockTestUser = entity.User{
		ID:             "test",
		FirstName:      "test",
		LastName:       "test",
		HashedPassword: "$2a$10$y01XUU2tPmc7.XpcB31v6e8GQPmkAODnCblk7LbYByuOuyyhptvaq", // hashed of "test"
		Email:          "test@test.com",
	}
)

func TestSessionImpl_Login(t *testing.T) {

	type args struct {
		givenEmail    string
		givenPassword string
		err           error
	}

	tcs := map[string]args{
		"success": {
			givenEmail:    mockTestUser.Email,
			givenPassword: "test",
		},
		"email not found": {
			givenEmail:    "nonexist@gmail.com",
			givenPassword: "test",
			err:           ErrEmailInvalid,
		},
		"invalid password": {
			givenEmail:    mockTestUser.Email,
			givenPassword: "wrongpassword",
			err:           ErrPasswordInvalid,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			database.TestTransaction(database.TestDB(), func(dbTx *database.DB) {
				s := NewSession(dbTx)
				dbTx.GormDB.Create(&mockTestUser)
				// When
				err := s.Login(tc.givenEmail, tc.givenPassword)
				// Then
				if tc.err != nil {
					require.Error(t, errors.Unwrap(err), tc.err.Error())
				} else {
					require.NoError(t, err)
					got := entity.User{}
					require.NoError(t, dbTx.GormDB.Find(&got).Error)
				}
			})
		})
	}
}
