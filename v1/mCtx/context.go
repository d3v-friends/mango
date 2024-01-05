package mCtx

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/d3v-friends/go-pure/fnPanic"
	"go.mongodb.org/mongo-driver/mongo"
)

const ctxMongo = "CTX_MONGO"

func SetDB(ctx context.Context, db *mongo.Database) context.Context {
	return context.WithValue(ctx, ctxMongo, db)
}

func GetDB(ctx context.Context) (db *mongo.Database, err error) {
	if ctx.Err() != nil {
		err = ctx.Err()
		return
	}

	var isOk bool
	if db, isOk = ctx.Value(ctxMongo).(*mongo.Database); !isOk {
		err = fmt.Errorf(
			"not found *mongo.Database in context.Context: ctx=%s",
			fnPanic.OnValue(json.Marshal(ctx)),
		)
		return
	}
	return
}

func GetDBP(ctx context.Context) *mongo.Database {
	return fnPanic.Get(GetDB(ctx))
}
