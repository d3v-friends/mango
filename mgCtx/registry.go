package mgCtx

import (
	"context"
	"github.com/d3v-friends/go-tools/fnCtx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
)

const ctxKeyMongoRegistry fnCtx.Key[*bsoncodec.Registry] = "CTX_MONGO_REGISTRY"

func SetRegistry(ctx context.Context, registry *bsoncodec.Registry) context.Context {
	return fnCtx.Set(ctx, ctxKeyMongoRegistry, registry)
}

func GetRegistry(ctx context.Context) (registry *bsoncodec.Registry) {
	var err error
	if registry, err = fnCtx.Get(ctx, ctxKeyMongoRegistry); err != nil {
		registry = bson.NewRegistry()
	}
	return
}
