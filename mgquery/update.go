package mgquery

import (
	"context"

	"github.com/d3v-friends/mango/v2"
	"github.com/d3v-friends/mango/v2/mgbuilder"
	"github.com/d3v-friends/mango/v2/mgctx"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func UpdateOne[T mango.Model](
	ctx context.Context,
	filter any,
	updater bson.M,
	opts ...*options.UpdateOneOptionsBuilder,
) (err error) {
	var col *mongo.Collection
	if col, err = mgctx.GetCol(ctx, *new(T)); err != nil {
		return
	}

	var f = mgbuilder.Filter(filter)

	var opt = options.UpdateOne()
	if len(opts) == 1 {
		opt = opts[0]
	}

	if _, err = col.UpdateOne(ctx, f, updater, opt); err != nil {
		return
	}

	return
}

func UpdateMany[T mango.Model](
	ctx context.Context,
	filter any,
	updater bson.M,
	opts ...*options.UpdateManyOptionsBuilder,
) (err error) {
	var col *mongo.Collection
	if col, err = mgctx.GetCol(ctx, *new(T)); err != nil {
		return
	}

	var f = mgbuilder.Filter(filter)

	var opt = options.UpdateMany()
	if len(opts) == 1 {
		opt = opts[0]
	}

	if _, err = col.UpdateMany(ctx, f, updater, opt); err != nil {
		return
	}

	return
}
