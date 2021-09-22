package database

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Transaction opens a transaction from db
func Transaction(db *DB, f func() error) error {
	return db.GormDB.Transaction(func(tx *gorm.DB) error {
		return errors.WithStack(f())
	})
}

// TestTransaction allows testing in db using transaction rollback
// First begin a transaction
// Run all mutation or query on database
// Assert result
// Rollback to ensure no actual commit is done
// Should not use in production
func TestTransaction(db DB, f func(dbTx *DB)) {
	// First begin a transaction
	tx := db.GormDB.Begin()

	// Wrap tx with db struct
	dbTx := &DB{GormDB: tx}

	// Run all mutation or query on database
	// Assert result
	f(dbTx)

	// Rollback to ensure no actual commit is done
	tx.Rollback()
}

// TestDB initializes and return db connection
// Should not use in production
func TestDB() DB {
	db, err := New("postgres://hydralink-api:hydralink-api@127.0.0.1:5432/hydralink-api?sslmode=disable")
	if err != nil {
		panic(err) // ideally should not use panic in production, but for testing is ok
	}
	return db
}
