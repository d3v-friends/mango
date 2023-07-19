package mcontext

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	CTX_MONGO_CLIENT = "CTX_MONGO_CLIENT"
)

func SetContext(ctx context.Context, client *mongo.Client) context.Context {
	return context.WithValue(ctx, CTX_MONGO_CLIENT, client)
}

func GetContext(ctx context.Context) *mongo.Client {
	cli, isOk := ctx.Value(CTX_MONGO_CLIENT).(*mongo.Client)
	if !isOk {
		panic(fmt.Errorf("not found mongo client in context"))
	}
	return cli
}
