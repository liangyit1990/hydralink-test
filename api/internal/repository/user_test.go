package repository

import (
	"testing"

	"github.com/hydralinkapp/hydralink/api/internal/entity"
	"github.com/hydralinkapp/hydralink/api/pkg/database"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	type args struct {
		given entity.User
		want  entity.User
	}

	tcs := map[string]args{
		"success": {
			given: entity.User{
				ID:             "test",
				FirstName:      "test",
				LastName:       "test",
				HashedPassword: "test",
				Email:          "test",
			},
			want: entity.User{
				ID:             "test",
				FirstName:      "test",
				LastName:       "test",
				HashedPassword: "test",
				Email:          "test",
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			database.TestTransaction(database.TestDB(), func(dbTx *database.DB) {
				// Given
				u := NewUser(dbTx)
				// When
				_, err := u.Create(tc.given)
				require.NoError(t, err)
				// Then
				got := &entity.User{}
				dbTx.GormDB.Where("id=?", "test").First(got)
				require.Equal(t, tc.want.ID, got.ID)
			})
		})
	}
}

func TestGetUser(t *testing.T) {
	type args struct {
		givenSearch UserSearch
		givenLock   bool
		want        entity.User
		err         error
	}

	mockUser := entity.User{
		ID:             "10000",
		FirstName:      "FirstName",
		LastName:       "LastName",
		HashedPassword: "HashedPassword",
		Email:          "Email",
	}

	tcs := map[string]args{
		"no result": {
			givenSearch: UserSearch{
				ID: "non-exist",
			},
		},
		"ok by id": {
			givenSearch: UserSearch{
				ID: mockUser.ID,
			},
			want: mockUser,
		},
		"ok by email": {
			givenSearch: UserSearch{
				Email: mockUser.Email,
			},
			want: mockUser,
		},
		"ok by first name": {
			givenSearch: UserSearch{
				FirstName: mockUser.FirstName,
			},
			want: mockUser,
		},
		"ok by last name": {
			givenSearch: UserSearch{
				LastName: mockUser.LastName,
			},
			want: mockUser,
		},
		"ok by all": {
			givenSearch: UserSearch{
				FirstName: mockUser.FirstName,
				LastName:  mockUser.LastName,
			},
			want: mockUser,
		},
		"ok by all - lock": {
			givenSearch: UserSearch{
				FirstName: mockUser.FirstName,
				LastName:  mockUser.LastName,
			},
			givenLock: true,
			want:      mockUser,
		},
	}

	for desc, tc := range tcs {
		tc := tc
		t.Run(desc, func(t *testing.T) {
			t.Parallel()
			database.TestTransaction(database.TestDB(), func(dbTx *database.DB) {
				// Given
				u := NewUser(dbTx)
				dbTx.GormDB.Create(&mockUser)
				// When
				got, err := u.Find(tc.givenSearch, tc.givenLock)
				// Then
				if tc.err != nil {
					require.Error(t, err, tc.err.Error())
				} else {
					require.NoError(t, err)
					require.Equal(t, tc.want.ID, got.ID)
					require.Equal(t, tc.want.FirstName, got.FirstName)
					require.Equal(t, tc.want.LastName, got.LastName)
					require.Equal(t, tc.want.HashedPassword, got.HashedPassword)
				}
			})
		})
	}
}

func TestUpdateUser(t *testing.T) {
	type args struct {
		givenNewPassword string
		want             entity.User
		err              error
	}

	mockUser := entity.User{
		ID:             "test",
		FirstName:      "FirstName",
		LastName:       "LastName",
		HashedPassword: "HashedPassword",
		Email:          "Email",
	}

	want := entity.User{
		FirstName:      "FirstName",
		LastName:       "LastName",
		HashedPassword: "Updated",
		Email:          "Email",
	}

	tcs := map[string]args{
		"success": {
			givenNewPassword: want.HashedPassword,
			want:             want,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			database.TestTransaction(database.TestDB(), func(dbTx *database.DB) {
				// Given
				u := NewUser(dbTx)
				dbTx.GormDB.Create(&mockUser)

				found := entity.User{}
				dbTx.GormDB.Find(&found)

				found.HashedPassword = tc.givenNewPassword

				// When
				err := u.Update(found)
				// Then
				if tc.err != nil {
					require.Error(t, err, tc.err.Error())
				} else {
					require.NoError(t, err)
					got := &entity.User{}
					require.NoError(t, dbTx.GormDB.First(got).Error)
					require.Equal(t, tc.want.FirstName, got.FirstName)
					require.Equal(t, tc.want.LastName, got.LastName)
					require.Equal(t, tc.want.HashedPassword, got.HashedPassword)
					require.Equal(t, tc.want.Email, got.Email)
				}
			})
		})
	}
}
