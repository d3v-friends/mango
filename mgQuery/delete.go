package mgQuery

import (
	"context"
	"github.com/d3v-friends/mango/mgCtx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DeleteOne[T Model](
	ctx context.Context,
	filter any,
	opts ...*options.DeleteOptions,
) (err error) {
	var col *mongo.Collection
	if col, err = mgCtx.GetColByModel[T](ctx); err != nil {
		return
	}

	var f bson.M
	if f, err = ParseFilter(filter); err != nil {
		return
	}

	if _, err = col.DeleteOne(ctx, f, opts...); err != nil {
		return
	}

	return
}

func DeleteMany[T Model](
	ctx context.Context,
	filter any,
	opts ...*options.DeleteOptions,
) (err error) {
	var col *mongo.Collection
	if col, err = mgCtx.GetColByModel[T](ctx); err != nil {
		return
	}

	var f bson.M
	if f, err = ParseFilter(filter); err != nil {
		return
	}

	if _, err = col.DeleteMany(ctx, f, opts...); err != nil {
		return
	}

	return
}
