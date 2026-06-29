package mgquery

import (
	"context"

	"github.com/d3v-friends/go-tools/fnError"
	"github.com/d3v-friends/mango/v2"
	"github.com/d3v-friends/mango/v2/mgctx"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	ErrEmptyModels = "empty_models"
)

// InsertOne
func InsertOne[T mango.Model](
	ctx context.Context,
	model T,
	opts ...*options.InsertOneOptionsBuilder,
) (_ *mongo.InsertOneResult, err error) {
	var col *mongo.Collection

	if col, err = mgctx.GetWriterCollection(ctx, model); err != nil {
		return
	}

	var opt = options.InsertOne()
	if len(opts) == 1 {
		opt = opts[0]
	}

	return col.InsertOne(ctx, model, opt)
}

func InsertMany[T mango.Model](
	ctx context.Context,
	models []T,
	opts ...*options.InsertManyOptionsBuilder,
) (_ *mongo.InsertManyResult, err error) {
	if len(models) == 0 {
		err = fnError.New(ErrEmptyModels)
		return
	}

	var ls = make([]interface{}, len(models))

	for i, model := range models {
		ls[i] = model
	}

	var col *mongo.Collection
	if col, err = mgctx.GetWriterCollection(ctx, models[0]); err != nil {
		return
	}

	var opt = options.InsertMany()
	if len(opts) == 1 {
		opt = opts[0]
	}

	return col.InsertMany(ctx, ls, opt)
}
