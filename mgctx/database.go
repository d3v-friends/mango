package mgctx

import (
	"context"

	"github.com/d3v-friends/go-tools/fnCtx"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetDB(ctx context.Context, db *mongo.Database) context.Context {
	return fnCtx.Set(ctx, ctxKeyMongoDB, db)
}

func GetDB(ctx context.Context) (*mongo.Database, error) {
	return fnCtx.Get(ctx, ctxKeyMongoDB)
}

func GetDBP(ctx context.Context) *mongo.Database {
	return fnCtx.GetP(ctx, ctxKeyMongoDB)
}

func GetClient(ctx context.Context) (_ *mongo.Client, err error) {
	var db *mongo.Database
	if db, err = fnCtx.Get(ctx, ctxKeyMongoDB); err != nil {
		return
	}
	return db.Client(), nil
}
