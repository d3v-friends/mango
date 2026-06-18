package mgquery

import (
	"context"

	"github.com/d3v-friends/go-tools/fnError"
	"github.com/d3v-friends/mango/v2"
	"github.com/d3v-friends/mango/v2/mgbuilder"
	"github.com/d3v-friends/mango/v2/mgctx"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func FindOne[T mango.Model](
	ctx context.Context,
	filter any,
	sorter any,
	opts ...*options.FindOneOptionsBuilder,
) (res *T, err error) {

	var f = mgbuilder.BsonM(filter)

	var o = options.FindOne()
	if len(opts) == 1 {
		o = opts[0]
	}

	var sort = mgbuilder.BsonD(sorter)

	if 0 < len(sort) {
		o.SetSort(sort)
	}

	var col *mongo.Collection
	if col, err = mgctx.GetCol(ctx, *new(T)); err != nil {
		return
	}

	var cur *mongo.SingleResult
	if cur = col.FindOne(ctx, f, o); cur.Err() != nil {
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
	opts ...*options.FindOptionsBuilder,
) (res []*T, err error) {

	var f = mgbuilder.BsonM(filter)

	var o = options.Find()
	if len(opts) == 1 {
		o = opts[0]
	}
	var sort = mgbuilder.BsonD(sorter)

	if 0 < len(sort) {
		o.SetSort(sort)
	}

	if limit != nil {
		o.SetLimit(*limit)
	}

	var col *mongo.Collection
	if col, err = mgctx.GetCol(ctx, *new(T)); err != nil {
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
	opts ...*options.FindOneAndUpdateOptionsBuilder,
) (res *T, err error) {
	var f = mgbuilder.BsonM(filter)

	var opt = options.FindOneAndUpdate()
	if len(opts) == 1 {
		opt = opts[0]
	}

	var sort = mgbuilder.BsonD(sorter)

	if 0 < len(sort) {
		opt.SetSort(sort)
	}

	var col *mongo.Collection
	if col, err = mgctx.GetCol(ctx, *new(T)); err != nil {
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
	pager mgbuilder.PagerArgs,
	opts ...*options.FindOptionsBuilder,
) (res *ModelList[T], err error) {
	var f = mgbuilder.BsonM(filter)

	var col *mongo.Collection
	if col, err = mgctx.GetCol(ctx, *new(T)); err != nil {
		return
	}

	var total int64
	if total, err = col.CountDocuments(ctx, f); err != nil {
		return
	}

	var o = options.Find()
	if len(opts) == 1 {
		o = opts[0]
	}

	var sort = mgbuilder.BsonD(sorter)

	if 0 < len(sort) {
		o.SetSort(sort)
	}

	if pager == nil {
		err = fnError.NewFields(ErrNotFoundPagerArgs, map[string]any{
			"filter": filter,
		})
		return
	}

	o.SetSkip(pager.GetSkip())
	o.SetLimit(pager.GetLimit())

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
