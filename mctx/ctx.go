package mctx

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

type IMongo interface {
	*mongo.Client | *mongo.Database
}

const CtxKey = "MONGO_DB"

func Set[T IMongo](ctx context.Context, v T) context.Context {
	return context.WithValue(ctx, CtxKey, v)
}

func Get[T IMongo](ctx context.Context) (res T, err error) {
	var has bool
	if res, has = ctx.Value(CtxKey).(T); !has {
		err = fmt.Errorf("not found mongo db in context")
		return
	}
	return
}

func GetP[T IMongo](ctx context.Context) (res T) {
	var err error
	if res, err = Get[T](ctx); err != nil {
		panic(err)
	}
	return
}
