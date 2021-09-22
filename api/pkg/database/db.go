package database

import (
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB wraps underlying database pkg
type DB struct {
	GormDB *gorm.DB
}

// New creates new DB
func New(connString string) (DB, error) {
	gormDB, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return DB{}, errors.WithStack(err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return DB{}, errors.WithStack(err)
	}

	if err := sqlDB.Ping(); err != nil {
		return DB{}, errors.WithStack(err)
	}

	// TODO : Configure this
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	return DB{GormDB: gormDB}, nil
}

// Close destroy any sql db connection
func (d DB) Close() error {
	sqlDB, err := d.GormDB.DB()
	if err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(sqlDB.Close())
}
