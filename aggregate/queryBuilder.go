package aggregate

import (
	"context"
	"fmt"
	"github.com/d3v-friends/mango/mtype"
	"github.com/d3v-friends/mango/mvars"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
)

type (
	QueryBuilder[MODEL mtype.IfColNm] struct {
		model    MODEL
		pipeline bson.A
	}

	FnQueryBuilder func(v any) (res bson.A, err error)
)

func NewQueryBuilder[MODEL mtype.IfColNm](
	fns ...FnQueryBuilder,
) (res *QueryBuilder[MODEL], err error) {
	res = &QueryBuilder[MODEL]{
		pipeline: make(bson.A, 0),
	}

	if err = res.validate(); err != nil {
		res = nil
		return
	}

	for _, fn := range fns {
		var query bson.A
		if query, err = fn(res.model); err != nil {
			return
		}
		res.pipeline = append(res.pipeline, query...)
	}

	return
}

func (x *QueryBuilder[MODEL]) FindOne(
	ctx context.Context,
	db *mongo.Database,
	filter any,
) (res *MODEL, err error) {
	var pipeline = make(bson.A, 0)
	pipeline = append(pipeline, bson.M{
		mvars.OMatch: filter,
	})
	pipeline = append(pipeline, x.pipeline...)

	var col = db.Collection(x.model.GetCollectionNm())
	var cursor *mongo.Cursor

	//fmt.Printf(
	//	"%s\n",
	//	fnPanic.OnValue(json.Marshal(pipeline)),
	//)

	if cursor, err = col.Aggregate(ctx, pipeline, &options.AggregateOptions{}); err != nil {
		return
	}

	var ls = make([]*MODEL, 0)
	if err = cursor.All(ctx, &ls); err != nil {
		return
	}

	if len(ls) == 0 {
		err = fmt.Errorf("not found data")
		return
	}

	res = ls[0]
	return
}

func (x *QueryBuilder[MODEL]) FindAll(
	ctx context.Context,
	db *mongo.Database,
	filter any,
) (res []*MODEL, err error) {
	var pipeline = make(bson.A, 0)
	pipeline = append(pipeline, bson.M{
		mvars.OMatch: filter,
	})
	pipeline = append(pipeline, x.pipeline...)

	var col = db.Collection(x.model.GetCollectionNm())
	var cursor *mongo.Cursor
	if cursor, err = col.Aggregate(ctx, pipeline, &options.AggregateOptions{}); err != nil {
		return
	}

	res = make([]*MODEL, 0)
	if err = cursor.All(ctx, &res); err != nil {
		return
	}

	return
}

func (x *QueryBuilder[MODEL]) FindList(
	ctx context.Context,
	db *mongo.Database,
	filter any,
	pager mtype.IfPager,
) (res *mtype.ResultList[MODEL], err error) {
	res = &mtype.ResultList[MODEL]{
		Page:  pager.GetPage(),
		Size:  pager.GetSize(),
		Total: 0,
		List:  make([]*MODEL, 0),
	}

	defer func() {
		if err != nil {
			res = nil
		}
	}()

	var pipeline = make(bson.A, 0)
	pipeline = append(pipeline, bson.M{
		mvars.OMatch: filter,
	})

	pipeline = append(pipeline, bson.M{
		mvars.OSkip: pager.GetSize() * pager.GetPage(),
	})

	pipeline = append(pipeline, bson.M{
		mvars.OLimit: pager.GetSize(),
	})

	pipeline = append(pipeline, x.pipeline...)

	var col = db.Collection(x.model.GetCollectionNm())
	var cursor *mongo.Cursor
	if cursor, err = col.Aggregate(ctx, pipeline, &options.AggregateOptions{}); err != nil {
		return
	}

	if err = cursor.All(ctx, &res.List); err != nil {
		return
	}

	if res.Total, err = col.CountDocuments(ctx, filter); err != nil {
		return
	}

	return
}

func (x *QueryBuilder[MODEL]) validate() (err error) {
	var model = new(MODEL)
	var to = reflect.TypeOf(model)
	if to.Kind() != reflect.Pointer {
		err = fmt.Errorf("invalid model: name=%s", to.Name())
		return
	}

	var vo = reflect.ValueOf(model)
	var elem = vo.Elem().Interface()

	var elemTo = reflect.TypeOf(elem)
	if elemTo.Kind() != reflect.Struct {
		err = fmt.Errorf("model kind is not struct: model.Kind()=%s", elemTo.Kind().String())
		return
	}

	x.model = *model

	return
}
