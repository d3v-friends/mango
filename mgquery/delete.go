package mgquery

import (
	"context"

	"github.com/d3v-friends/mango/v2"
	"github.com/d3v-friends/mango/v2/mgbuilder"
	"github.com/d3v-friends/mango/v2/mgctx"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func DeleteOne[T mango.Model](
	ctx context.Context,
	filter any,
	opts ...*options.DeleteOneOptionsBuilder,
) (res *mongo.DeleteResult, err error) {
	var col *mongo.Collection
	if col, err = mgctx.GetWriterCollection(ctx, *new(T)); err != nil {
		return
	}

	var f = mgbuilder.Filter(ctx, filter)

	var opt = options.DeleteOne()
	if len(opts) == 1 {
		opt = opts[0]
	}

	if res, err = col.DeleteOne(ctx, f, opt); err != nil {
		return
	}

	return
}

func DeleteMany[T mango.Model](
	ctx context.Context,
	filter any,
	opts ...*options.DeleteManyOptionsBuilder,
) (res *mongo.DeleteResult, err error) {
	var col *mongo.Collection
	if col, err = mgctx.GetWriterCollection(ctx, *new(T)); err != nil {
		return
	}

	var f = mgbuilder.Filter(ctx, filter)

	var opt = options.DeleteMany()
	if len(opts) == 1 {
		opt = opts[0]
	}

	if res, err = col.DeleteMany(ctx, f, opt); err != nil {
		return
	}

	return
}
