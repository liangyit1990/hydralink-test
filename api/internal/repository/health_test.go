package repository

import (
	"testing"

	"github.com/hydralinkapp/hydralink/api/pkg/database"
	"github.com/stretchr/testify/require"
)

func TestHealth_Readiness(t *testing.T) {
	t.Run("readiness", func(t *testing.T) {
		db := database.TestDB()
		h := NewHealth(&db)
		require.NoError(t, h.Readiness())
	})
}
