package fnMango

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

const defClientKey = "MONGO_CLIENT"
const defDbKey = "MONGO_DB"

func GetClient(ctx context.Context, ctxKeys ...string) (client *mongo.Client, err error) {
	var key = defClientKey
	if 0 < len(ctxKeys) {
		key = ctxKeys[0]
	}

	var has bool
	if client, has = ctx.Value(key).(*mongo.Client); !has {
		err = fmt.Errorf("not found *mongo.Client")
		return
	}

	return
}

func GetClientP(ctx context.Context, ctxKeys ...string) (client *mongo.Client) {
	var err error
	if client, err = GetClient(ctx, ctxKeys...); err != nil {
		panic(err)
	}
	return
}

func SetClient(ctx context.Context, client *mongo.Client, ctxKeys ...string) context.Context {
	var key = defClientKey
	if 0 < len(ctxKeys) {
		key = ctxKeys[0]
	}
	return context.WithValue(ctx, key, client)
}

func GetDb(ctx context.Context, ctxKeys ...string) (db *mongo.Database, err error) {
	var key = defClientKey
	if 0 < len(ctxKeys) {
		key = ctxKeys[0]
	}

	var has bool
	if db, has = ctx.Value(key).(*mongo.Database); !has {
		err = fmt.Errorf("not found *mongo.Database")
		return
	}

	return
}

func GetDbP(ctx context.Context, ctxKeys ...string) (db *mongo.Database) {
	var err error
	if db, err = GetDb(ctx, ctxKeys...); err != nil {
		panic(err)
	}
	return
}

func SetDb(ctx context.Context, db *mongo.Database, ctxKeys ...string) context.Context {
	var key = defClientKey
	if 0 < len(ctxKeys) {
		key = ctxKeys[0]
	}
	return context.WithValue(ctx, key, db)
}
