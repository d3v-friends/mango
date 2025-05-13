package mgQuery

import (
	"context"
	"github.com/d3v-friends/go-tools/fnError"
	"github.com/d3v-friends/go-tools/fnPointer"
	"github.com/d3v-friends/mango"
	"github.com/d3v-friends/mango/mgCtx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindOne[T mango.Model](
	ctx context.Context,
	filter any,
	sorter any,
	opts ...*options.FindOneOptions,
) (res *T, err error) {
	var f any
	if f, err = ParseFilter(filter); err != nil {
		return
	}

	var opt = &options.FindOneOptions{}
	if len(opts) == 1 {
		opt = opts[0]
	}

	if opt.Sort, err = ParseSorter(sorter); err != nil {
		return
	}

	var col *mongo.Collection
	if col, err = mgCtx.GetCol(ctx, *new(T)); err != nil {
		return
	}

	var cur *mongo.SingleResult
	if cur = col.FindOne(ctx, f, opts...); cur.Err() != nil {
		err = cur.Err()
		return
	}

	res = new(T)
	if err = cur.Decode(res); err != nil {
		return
	}

	return
}

func Find[T mango.Model](
	ctx context.Context,
	filter any,
	sorter any,
	limit *int64,
	opts ...*options.FindOptions,
) (res []*T, err error) {
	var f any
	if f, err = ParseFilter(filter); err != nil {
		return
	}

	var o = &options.FindOptions{}
	if len(opts) == 1 {
		o = opts[0]
	}

	if o.Sort, err = ParseSorter(sorter); err != nil {
		return
	}

	o.Limit = limit

	var col *mongo.Collection
	if col, err = mgCtx.GetCol(ctx, *new(T)); err != nil {
		return
	}

	var cur *mongo.Cursor
	if cur, err = col.Find(ctx, f, o); err != nil {
		return
	}

	res = make([]*T, 0)
	if err = cur.All(ctx, &res); err != nil {
		return
	}

	return
}

func FindOneAndUpdate[T mango.Model](
	ctx context.Context,
	filter any,
	sorter any,
	updater bson.M,
	opts ...*options.FindOneAndUpdateOptions,
) (res *T, err error) {
	var f any
	if f, err = ParseFilter(filter); err != nil {
		return
	}

	var opt = &options.FindOneAndUpdateOptions{}
	if len(opts) == 1 {
		opt = opts[0]
	}

	if opt.Sort, err = ParseSorter(sorter); err != nil {
		return
	}

	var col *mongo.Collection
	if col, err = mgCtx.GetCol(ctx, *new(T)); err != nil {
		return
	}

	var cur *mongo.SingleResult
	if cur = col.
		FindOneAndUpdate(ctx, f, updater, opt); cur.Err() != nil {
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

func FindList[T mango.Model](
	ctx context.Context,
	filter any,
	sorter any,
	pager PagerArgs,
	opts ...*options.FindOptions,
) (res *ModelList[T], err error) {
	var f any
	if f, err = ParseFilter(filter); err != nil {
		return
	}

	var col *mongo.Collection
	if col, err = mgCtx.GetCol(ctx, *new(T)); err != nil {
		return
	}

	var total int64
	if total, err = col.CountDocuments(ctx, f); err != nil {
		return
	}

	var o = &options.FindOptions{}
	if len(opts) == 1 {
		o = opts[0]
	}

	if o.Sort, err = ParseSorter(sorter); err != nil {
		return
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
	if cur, err = col.Find(ctx, f, o); err != nil {
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
