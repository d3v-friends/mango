package mango

import (
	"context"
	"github.com/d3v-friends/go-pure/fnParams"
	"github.com/d3v-friends/mango/mFilter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateOne[MODEL any](
	ctx context.Context,
	i IfFilter,
	u IfUpdate,
	opts ...*options.UpdateOptions,
) (_ *MODEL, err error) {
	var filter bson.M
	if filter, err = i.Filter(); err != nil {
		return
	}

	var update bson.M
	if update, err = u.Update(); err != nil {
		return
	}

	var col = GetMangoP(ctx).DB.Collection(i.ColNm())
	if _, err = col.UpdateOne(
		ctx,
		filter,
		update,
		fnParams.Get(opts),
	); err != nil {
		return
	}

	return ReadOne[MODEL](ctx, i)
}

func UpdateMany[MODEL any](
	ctx context.Context,
	i IfFilter,
	u IfUpdate,
	opts ...*options.UpdateOptions,
) (ls []*MODEL, err error) {
	var filter bson.M
	if filter, err = i.Filter(); err != nil {
		return
	}

	var update bson.M
	if update, err = u.Update(); err != nil {
		return
	}

	var col = GetMangoP(ctx).DB.Collection(i.ColNm())
	var res *mongo.UpdateResult
	if res, err = col.UpdateMany(
		ctx,
		filter,
		update,
		fnParams.Get(opts),
	); err != nil {
		return
	}

	var updateFilter IfFilter
	switch id := res.UpsertedID.(type) {
	case primitive.ObjectID:
		updateFilter = mFilter.NewId(i.ColNm(), id)
	case []primitive.ObjectID:
		updateFilter = mFilter.NewIds(i.ColNm(), mFilter.OperatorIn, id...)
	default:
		ls = make([]*MODEL, 0)
		return
	}

	return ReadAll[MODEL](ctx, updateFilter)
}
