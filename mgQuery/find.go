package mgQuery

import (
	"context"
	"github.com/d3v-friends/go-tools/fnError"
	"github.com/d3v-friends/go-tools/fnPointer"
	"github.com/d3v-friends/mango/mgCtx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindOne[T Model](
	ctx context.Context,
	filter bson.M,
	opts ...*options.FindOneOptions,
) (res *T, err error) {
	var col *mongo.Collection
	if col, err = mgCtx.GetColByModel[T](ctx); err != nil {
		return
	}

	var cur *mongo.SingleResult
	if cur = col.FindOne(ctx, filter, opts...); cur.Err() != nil {
		err = cur.Err()
		return
	}

	res = new(T)
	if err = cur.Decode(res); err != nil {
		return
	}

	return
}

func Find[T Model](
	ctx context.Context,
	filter bson.M,
	limit *int64,
	opts ...*options.FindOptions,
) (res []*T, err error) {
	var o = &options.FindOptions{}
	if len(opts) == 1 {
		o = opts[0]
	}

	o.Limit = limit

	var col *mongo.Collection
	if col, err = mgCtx.GetColByModel[T](ctx); err != nil {
		return
	}

	var cur *mongo.Cursor
	if cur, err = col.Find(ctx, filter, o); err != nil {
		return
	}

	res = make([]*T, 0)
	if err = cur.All(ctx, &res); err != nil {
		return
	}

	return
}

func FindOneAndUpdate[T Model](
	ctx context.Context,
	filter bson.M,
	updater bson.M,
	opts ...*options.FindOneAndUpdateOptions,
) (res *T, err error) {
	var opt = &options.FindOneAndUpdateOptions{}
	if len(opts) == 1 {
		opt = opts[0]
	} else {
		opt.ReturnDocument = fnPointer.Make(options.After)
	}

	var col *mongo.Collection
	if col, err = mgCtx.GetColByModel[T](ctx); err != nil {
		return
	}

	var cur *mongo.SingleResult
	if cur = col.
		FindOneAndUpdate(ctx, filter, updater, opt); cur.Err() != nil {
		err = cur.Err()
		return
	}

	res = new(T)
	if err = cur.Decode(res); err != nil {
		return
	}

	return
}

type ModelList[T any] struct {
	Page  int64
	Size  int64
	Total int64
	List  []*T
}

const ErrNotFoundPagerArgs = "not_found_pager_args"

func FindList[T Model](
	ctx context.Context,
	filter bson.M,
	pager PagerArgs,
	opts ...*options.FindOptions,
) (res *ModelList[T], err error) {
	var col *mongo.Collection
	if col, err = mgCtx.GetColByModel[T](ctx); err != nil {
		return
	}

	var total int64
	if total, err = col.CountDocuments(ctx, filter); err != nil {
		return
	}

	var o = &options.FindOptions{}
	if len(opts) == 1 {
		o = opts[0]
	}

	if fnPointer.IsNil(pager) {
		err = fnError.NewFields(ErrNotFoundPagerArgs, map[string]any{
			"filter": filter,
		})
		return
	}

	o.Skip = fnPointer.Make(pager.GetSize() * pager.GetPage())
	o.Limit = fnPointer.Make(pager.GetSize())

	var cur *mongo.Cursor
	if cur, err = col.Find(ctx, filter, o); err != nil {
		return
	}

	var list = make([]*T, 0)
	if err = cur.All(ctx, &list); err != nil {
		return
	}

	res = &ModelList[T]{
		Page:  pager.GetPage(),
		Size:  pager.GetSize(),
		Total: total,
		List:  list,
	}

	return
}
