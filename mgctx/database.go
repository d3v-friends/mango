package mgctx

import (
	"context"

	"github.com/d3v-friends/go-tools/fnCtx"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetDatabase(ctx context.Context, db *mongo.Database) context.Context {
	return fnCtx.Set(ctx, ctxKeyMongoDB, db)
}

func GetDatabase(ctx context.Context) (*mongo.Database, error) {
	return fnCtx.Get(ctx, ctxKeyMongoDB)
}

func GetReaderDB(ctx context.Context) (db *mongo.Database, err error) {
	if db, err = fnCtx.Get(ctx, ctxKeyReaderMongoDB); err == nil {
		return
	}
	return fnCtx.Get(ctx, ctxKeyMongoDB)
}

func GetWriterDB(ctx context.Context) (db *mongo.Database, err error) {
	if db, err = fnCtx.Get(ctx, ctxKeyWriterMongoDB); err == nil {
		return
	}
	return fnCtx.Get(ctx, ctxKeyMongoDB)
}

func SetReaderDB(ctx context.Context, db *mongo.Database) context.Context {
	return fnCtx.Set(ctx, ctxKeyReaderMongoDB, db)
}

func SetWriterDB(ctx context.Context, db *mongo.Database) context.Context {
	return fnCtx.Set(ctx, ctxKeyWriterMongoDB, db)
}
