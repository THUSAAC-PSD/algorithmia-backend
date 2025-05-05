package contract

import "context"

type UnitOfWork interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit() error
	Rollback() error
}

type UnitOfWorkFactory interface {
	New() UnitOfWork
}
