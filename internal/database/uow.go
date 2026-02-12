package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type txKey struct{}

type UnitOfWork interface {
	Commit() error
	Rollback() error
}

type UnitOfWorkFactory interface {
	Begin(context.Context) (UnitOfWork, context.Context, error)
}

type SqlxUnitOfWorkFactory struct {
	db *sqlx.DB
}

func NewUnitOfWorkFactory(db *sqlx.DB) *SqlxUnitOfWorkFactory {
	return &SqlxUnitOfWorkFactory{db: db}
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

type SqlxUnitOfWork struct {
	tx *sqlx.Tx
}

func (u *SqlxUnitOfWork) Commit() error {
	return u.tx.Commit()
}

func (u *SqlxUnitOfWork) Rollback() error {
	return u.tx.Rollback()
}

func GetDB(ctx context.Context, factory *SqlxUnitOfWorkFactory) any {
	if tx, ok := ctx.Value(txKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return factory.GetDB()
}

func (f *SqlxUnitOfWorkFactory) GetDB() *sqlx.DB {
	return f.db
}
