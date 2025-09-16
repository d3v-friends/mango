package mgQuery

import (
	"context"

	"github.com/d3v-friends/mango"
	"github.com/d3v-friends/mango/mgCtx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Count[T mango.Model](
	ctx context.Context,
	filter any,
	opts ...*options.CountOptions,
) (_ int64, err error) {
	var registry = mgCtx.GetRegistry(ctx)

	var f any
	if f, err = ParseFilter(filter, registry); err != nil {
		return
	}

	var col *mongo.Collection
	if col, err = mgCtx.GetCol(ctx, *new(T)); err != nil {
		return
	}

	return col.CountDocuments(ctx, f, opts...)
}
