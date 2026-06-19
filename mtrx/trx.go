package mtrx

import (
	"context"

	"github.com/d3v-friends/go-tools/fnError"
	"github.com/d3v-friends/mango/v2/mgctx"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	ErrInvalidResultType = "invalid_result_type"
)

type FnTrx[T any] = func(ctx context.Context) (res T, err error)

func Do[T any](
	ctx context.Context,
	fn FnTrx[T],
	opts ...*options.SessionOptionsBuilder,
) (_ T, err error) {
	var db *mongo.Database
	if db, err = mgctx.GetDB(ctx); err != nil {
		return
	}

	var opt = options.Session()
	if len(opts) == 1 {
		opt = opts[0]
	}

	var session *mongo.Session
	if session, err = db.Client().StartSession(opt); err != nil {
		return
	}

	var res any
	if res, err = session.WithTransaction(ctx, func(ctx context.Context) (interface{}, error) {
		return fn(ctx)
	}); err != nil {
		return
	}

	var r, ok = res.(T)
	if !ok {
		err = fnError.NewF(ErrInvalidResultType)
		return
	}

	return r, nil
}
