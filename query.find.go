package mango

import (
	"context"
	"github.com/d3v-friends/go-pure/fnParams"
	"github.com/d3v-friends/go-pure/fnReflect"
	"go.mongodb.org/mongo-driver/bson"
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
