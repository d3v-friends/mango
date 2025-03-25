package mgQuery

import (
	"context"
	"github.com/d3v-friends/go-tools/fnError"
	"github.com/d3v-friends/mango/mgCtx"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ErrEmptyModels                    = "empty_models"
	ErrModelsAreNotSameCollectionData = "models_are_not_same_collection_data"
)

func InsertOne[T Model](
	ctx context.Context,
	model T,
) (err error) {
	var col *mongo.Collection

	if col, err = mgCtx.GetColByModel[T](ctx); err != nil {
		return
	}

	if _, err = col.InsertOne(ctx, model); err != nil {
		return
	}

	return
}

func InsertMany[T Model](
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
	if col, err = mgCtx.GetColByModel[T](ctx); err != nil {
		return
	}

	if _, err = col.InsertMany(ctx, ls); err != nil {
		return
	}

	return
}
