package mango

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/d3v-friends/go-pure/fnPanic"
	"github.com/d3v-friends/go-pure/fnParams"
	"github.com/d3v-friends/go-pure/fnReflect"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ReadOne[MODEL any](
	ctx context.Context,
	i IfFilter,
	opts ...*options.FindOneOptions,
) (res *MODEL, err error) {
	var filter bson.M
	if filter, err = i.Filter(); err != nil {
		return
	}

	var col = GetMangoP(ctx).DB.Collection(i.ColNm())
	var single *mongo.SingleResult
	if single = col.FindOne(
		ctx,
		filter,
		fnParams.Get(opts),
	); single.Err() != nil {
		err = single.Err()
		return
	}

	res = new(MODEL)
	if err = single.Decode(res); err != nil {
		return
	}

	return
}

func ReadAll[MODEL any](
	ctx context.Context,
	i IfFilter,
	opts ...*options.FindOptions,
) (ls []*MODEL, err error) {
	var filter bson.M
	if filter, err = i.Filter(); err != nil {
		return
	}

	var col = GetMangoP(ctx).DB.Collection(i.ColNm())
	var cur *mongo.Cursor
	if cur, err = col.Find(
		ctx,
		filter,
		fnParams.Get(opts),
	); err != nil {
		return
	}

	ls = make([]*MODEL, 0)
	if err = cur.All(ctx, &ls); err != nil {
		return
	}

	return
}

func ReadList[MODEL any](
	ctx context.Context,
	i IfFilter,
	p IfPager,
	opts ...*options.FindOptions,
) (ls []*MODEL, total int64, err error) {
	var filter bson.M
	if filter, err = i.Filter(); err != nil {
		return
	}

	var col = GetMangoP(ctx).DB.Collection(i.ColNm())
	if total, err = col.CountDocuments(ctx, filter); err != nil {
		return
	}

	var opt = fnParams.Get(opts)
	if opt == nil {
		opt = &options.FindOptions{}
	}

	opt.Skip = fnReflect.ToPointer(p.Page() * p.Size())
	opt.Limit = fnReflect.ToPointer(p.Size())

	var cur *mongo.Cursor
	if cur, err = col.Find(
		ctx,
		filter,
		opt,
	); err != nil {
		return
	}

	ls = make([]*MODEL, 0)
	if err = cur.All(ctx, &ls); err != nil {
		return
	}

	return
}

/* ------------------------------------------------------------------------------------------------------------ */

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
		updateFilter = NewIdFilter(i.ColNm(), id)
	case []primitive.ObjectID:
		updateFilter = NewIdsFilter(i.ColNm(), OperatorIn, id...)
	default:
		ls = make([]*MODEL, 0)
		return
	}

	return ReadAll[MODEL](ctx, updateFilter)
}

/* ------------------------------------------------------------------------------------------------------------ */

func DeleteOne(ctx context.Context, i IfFilter, opts ...*options.DeleteOptions) (err error) {
	var filter bson.M
	if filter, err = i.Filter(); err != nil {
		return
	}

	var col = GetMangoP(ctx).DB.Collection(i.ColNm())
	var count int64
	if count, err = col.CountDocuments(ctx, filter); err != nil {
		return
	}

	if count == 0 || 2 < count {
		err = fmt.Errorf("delete one must has 1 document: doc_count=%d, filter=%s", count, fnPanic.Get(json.Marshal(filter)))
		return
	}

	if _, err = col.DeleteOne(ctx, i, opts...); err != nil {
		return
	}

	return
}

func DeleteMany(ctx context.Context, i IfFilter, opts ...*options.DeleteOptions) (err error) {
	var filter bson.M
	if filter, err = i.Filter(); err != nil {
		return
	}

	var col = GetMangoP(ctx).DB.Collection(i.ColNm())
	if _, err = col.DeleteMany(ctx, filter, opts...); err != nil {
		return
	}

	return
}
