package mgQuery

import (
	"context"
	"github.com/d3v-friends/mango/mgCtx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateOne[T Model](
	ctx context.Context,
	filter any,
	updater bson.M,
	opts ...*options.UpdateOptions,
) (err error) {
	var col *mongo.Collection
	if col, err = mgCtx.GetColByModel[T](ctx); err != nil {
		return
	}

	var f bson.M
	if f, err = ParseFilter(filter); err != nil {
		return
	}

	if _, err = col.UpdateOne(ctx, f, updater, opts...); err != nil {
		return
	}

	return
}

func UpdateMany[T Model](
	ctx context.Context,
	filter any,
	updater bson.M,
	opts ...*options.UpdateOptions,
) (err error) {
	var col *mongo.Collection
	if col, err = mgCtx.GetColByModel[T](ctx); err != nil {
		return
	}

	var f bson.M
	if f, err = ParseFilter(filter); err != nil {
		return
	}

	if _, err = col.UpdateMany(ctx, f, updater, opts...); err != nil {
		return
	}

	return
}
