package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type txKey struct{}

// UnitOfWorkFactory is the factory/creator interface
// begin is how a transaction/unit of work is created
type UnitOfWorkFactory interface {
	Begin(context.Context) (UnitOfWork, context.Context, error)
	GetDB(context.Context) any
}

// UnitOfWork is the product interface
type UnitOfWork interface {
	Commit() error
	Rollback() error
}

// SqlxUnitOfWorkFactory is the concrete factory/creator
type SqlxUnitOfWorkFactory struct {
	db *sqlx.DB
}

func (f *SqlxUnitOfWorkFactory) GetDB(ctx context.Context) any {
	if tx, ok := ctx.Value(txKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return f.db
}

func (f *SqlxUnitOfWorkFactory) Begin(ctx context.Context) (UnitOfWork, context.Context, error) {
	tx, err := f.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, ctx, err
	}

	uow := &SqlxUnitOfWork{tx: tx}
	ctx = context.WithValue(ctx, txKey{}, tx)

	return uow, ctx, nil
}

func NewUnitOfWorkFactory(db *sqlx.DB) *SqlxUnitOfWorkFactory {
	return &SqlxUnitOfWorkFactory{db: db}
}

// SqlxUnitOfWork is the concrete product
type SqlxUnitOfWork struct {
	tx *sqlx.Tx
}

func (u *SqlxUnitOfWork) Commit() error {
	return u.tx.Commit()
}

func (u *SqlxUnitOfWork) Rollback() error {
	return u.tx.Rollback()
}
