package gorms

import (
	"gorm.io/gorm"
)

type Transactor interface {
	Action(func(conn DbConn) error) error
}

type Transaction struct {
	conn DbConn
}

func (t *Transaction) Action(f func(conn DbConn) error) error {
	t.conn.Begin()
	err := f(t.conn)
	if err != nil {
		t.conn.Rollback()
		return err
	}
	t.conn.Commit()
	return nil
}

func NewTransaction(db *gorm.DB) *Transaction {
	return &Transaction{
		conn: New(db),
	}
}
