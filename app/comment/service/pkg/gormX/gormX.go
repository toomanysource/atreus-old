package gormX

import (
	"context"
	"gorm.io/gorm"
)

type Transactor = *gorm.DB

type DB struct {
	db *gorm.DB
}

func NewConn(db *gorm.DB) *DB {
	return &DB{db: db}
}

func (c *DB) Action(ctx context.Context, f func(Transactor) error) error {
	err := c.db.Begin().Error
	if err != nil {
		return err
	}
	err = f(c.db.WithContext(ctx))
	if err != nil {
		c.db.Rollback()
		return err
	}
	err = c.db.Commit().Error
	if err != nil {
		c.db.Rollback()
		return err
	}
	return nil
}
