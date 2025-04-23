package mgCtx

import (
	"context"
	"github.com/d3v-friends/go-tools/fnCtx"
	"go.mongodb.org/mongo-driver/mongo"
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
