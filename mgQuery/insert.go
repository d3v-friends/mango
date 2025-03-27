package mgQuery

import (
	"context"
	"github.com/d3v-friends/go-tools/fnError"
	"github.com/d3v-friends/mango"
	"github.com/d3v-friends/mango/mgCtx"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ErrEmptyModels = "empty_models"
)

func InsertOne[T mango.Model](
	ctx context.Context,
	model T,
) (err error) {
	var col *mongo.Collection

	if col, err = mgCtx.GetCol(ctx, *new(T)); err != nil {
		return
	}

	if _, err = col.InsertOne(ctx, model); err != nil {
		return
	}

	return
}

func InsertMany[T mango.Model](
	ctx context.Context,
	models []T,
) (err error) {
	if len(models) == 0 {
		err = fnError.New(ErrEmptyModels)
		return
	}

	var ls = make([]interface{}, len(models))

	for i, model := range models {
		ls[i] = model
	}

	var col *mongo.Collection
	if col, err = mgCtx.GetCol(ctx, *new(T)); err != nil {
		return
	}

	if _, err = col.InsertMany(ctx, ls); err != nil {
		return
	}

	return
}
