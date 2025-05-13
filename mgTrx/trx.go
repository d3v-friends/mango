package mgTrx

import (
	"context"
	"github.com/d3v-friends/go-tools/fnError"
	"github.com/d3v-friends/mango/mgCtx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ErrInvalidResultType = "invalid_result_type"
)

func Do[T any](
	ctx context.Context,
	fn func(sessCtx mongo.SessionContext) (T, error),
	opts ...*options.SessionOptions,
) (_ T, err error) {
	var db *mongo.Database
	if db, err = mgCtx.GetDB(ctx); err != nil {
		return
	}

	var session mongo.Session
	if session, err = db.Client().StartSession(opts...); err != nil {
		return
	}

	if err = session.StartTransaction(); err != nil {
		return
	}

	defer session.EndSession(ctx)

	var res any
	if res, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return fn(sessCtx)
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
