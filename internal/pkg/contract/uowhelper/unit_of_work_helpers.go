package uowhelper

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"emperror.dev/errors"
)

func Do(
	ctx context.Context,
	uow contract.UnitOfWork,
	l logger.Logger,
	fn func(innerCtx context.Context) error,
) error {
	_, err := DoWithResult[int](ctx, uow, l, func(innerCtx context.Context) (int, error) {
		return 0, fn(innerCtx)
	})

	return err
}

func DoWithResult[T any](
	ctx context.Context,
	uow contract.UnitOfWork,
	l logger.Logger,
	fn func(innerCtx context.Context) (T, error),
) (T, error) {
	uowCtx, err := uow.Begin(ctx)
	if err != nil {
		var emptyT T
		return emptyT, errors.WrapIf(err, "failed to begin unit of work")
	}

	defer func() {
		if r := recover(); r != nil {
			if rbErr := uow.Rollback(); rbErr != nil {
				l.Error("failed to rollback unit of work after panic", errors.WithStack(rbErr))
			}

			panic(r)
		}
	}()

	res, err := fn(uowCtx)
	if err == nil {
		if commitErr := uow.Commit(); commitErr != nil {
			err = errors.WrapIf(err, "failed to commit unit of work")
		}

		return res, err
	}

	if rbErr := uow.Rollback(); rbErr != nil {
		l.Error("failed to rollback unit of work", errors.WithStack(rbErr))
	}

	return res, err
}
