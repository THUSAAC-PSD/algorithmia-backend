package database

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"emperror.dev/errors"
	"gorm.io/gorm"
)

type GormUnitOfWorkFactory struct {
	db *gorm.DB
	l  logger.Logger
}

func NewGormUnitOfWorkFactory(db *gorm.DB, l logger.Logger) contract.UnitOfWorkFactory {
	return &GormUnitOfWorkFactory{
		db: db,
		l:  l,
	}
}

func (g *GormUnitOfWorkFactory) New() contract.UnitOfWork {
	return NewGormUnitOfWork(g.db, g.l)
}

type txContextKey struct{}

type gormUnitOfWork struct {
	db *gorm.DB
	tx *gorm.DB
	l  logger.Logger
}

func NewGormUnitOfWork(db *gorm.DB, l logger.Logger) contract.UnitOfWork {
	return &gormUnitOfWork{
		db: db,
		l:  l,
	}
}

func (uow *gormUnitOfWork) Begin(ctx context.Context) (context.Context, error) {
	if uow.tx != nil {
		return nil, errors.New("transaction already started for this UoW instance")
	}

	tx := uow.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, errors.WrapIf(tx.Error, "failed to begin transaction")
	}

	uow.tx = tx

	txCtx := context.WithValue(ctx, txContextKey{}, uow.tx)
	return txCtx, nil
}

func (uow *gormUnitOfWork) Commit() error {
	if uow.tx == nil {
		return errors.New("no transaction to commit")
	}

	tx := uow.tx
	uow.tx = nil

	err := tx.Commit().Error
	return errors.WrapIf(err, "failed to commit transaction")
}

func (uow *gormUnitOfWork) Rollback() error {
	if uow.tx == nil {
		return nil
	}

	tx := uow.tx
	uow.tx = nil

	err := tx.Rollback().Error
	return errors.WrapIf(err, "failed to rollback transaction")
}

func GetDBFromContext(ctx context.Context, baseDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txContextKey{}).(*gorm.DB); ok && tx != nil {
		return tx.WithContext(ctx)
	}

	return baseDB.WithContext(ctx)
}
