package mango

import (
	"context"
	"fmt"
	"github.com/d3v-friends/go-pure/fnLogger"
	"github.com/d3v-friends/go-pure/fnParams"
	"github.com/d3v-friends/go-pure/fnReflect"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type (
	MDoc[DATA any] struct {
		Id        primitive.ObjectID   `bson:"_id"`
		Data      *DATA                `bson:"data"`
		History   []*MDocHistory[DATA] `bson:"history"`
		CreatedAt time.Time            `bson:"createdAt"`
		UpdatedAt time.Time            `bson:"updatedAt"`

		// managed
		origin *MDoc[DATA] `bson:"-"`
		colNm  string      `bson:"colNm"`
	}

	MDocHistory[DATA any] struct {
		Data      *DATA     `bson:"data"`
		CreatedAt time.Time `bson:"createdAt"`
	}
)

func (x *MDoc[DATA]) Update(ctx context.Context, v *DATA) (err error) {
	if x.origin == nil {
		err = fmt.Errorf(
			"this model is not loaded data: model=%s",
			fnLogger.ToJsonP(x),
		)
		return
	}

	var now = time.Now()
	var history = &MDocHistory[DATA]{
		Data:      x.Data,
		CreatedAt: x.UpdatedAt,
	}
	var update *mongo.UpdateResult
	if update, err = GetColP(ctx, x.colNm).UpdateByID(ctx, x.Id, bson.M{
		"$set": bson.M{
			"data":      v,
			"updatedAt": now,
		},
		"$push": bson.M{
			"history": history,
		},
	}); err != nil {
		return
	}

	var id, isOk = update.UpsertedID.(primitive.ObjectID)
	if !isOk || id != x.Id {
		err = fmt.Errorf("fail update model: model=%s", fnLogger.ToJsonP(x))
		return
	}

	// update this data
	x.Data = v
	x.History = append(x.History, history)
	x.UpdatedAt = now

	return
}

/* ------------------------------------------------------------------------------------------------------------ */

func ReadOneM[DATA any](
	ctx context.Context,
	i IfFilter,
	opts ...*options.FindOneOptions,
) (res *MDoc[DATA], err error) {
	var filter bson.M
	if filter, err = i.Filter(); err != nil {
		return
	}

	var sres *mongo.SingleResult
	if sres = GetColP(ctx, i.ColNm()).
		FindOne(ctx, filter, opts...); sres.Err() != nil {
		err = sres.Err()
		return
	}

	res = new(MDoc[DATA])
	if err = sres.Decode(res); err != nil {
		return
	}

	var origin = *res
	res.origin = &origin
	res.colNm = i.ColNm()

	return
}

func ReadAllM[DATA any](
	ctx context.Context,
	i IfFilter,
	opts ...*options.FindOptions,
) (ls []*MDoc[DATA], err error) {
	var filter bson.M
	if filter, err = i.Filter(); err != nil {
		return
	}

	var cur *mongo.Cursor
	if cur, err = GetColP(ctx, i.ColNm()).Find(ctx, filter, opts...); err != nil {
		return
	}

	ls = make([]*MDoc[DATA], 0)
	if err = cur.All(ctx, &ls); err != nil {
		return
	}

	for _, l := range ls {
		var origin = *l
		l.origin = &origin
		l.colNm = i.ColNm()
	}

	return
}

func ReadListM[DATA any](
	ctx context.Context,
	i IfFilter,
	p IfPager,
	opts ...*options.FindOptions,
) (
	ls []*MDoc[DATA],
	total int64,
	err error,
) {
	var filter bson.M
	if filter, err = i.Filter(); err != nil {
		return
	}

	if total, err = GetColP(ctx, i.ColNm()).
		CountDocuments(ctx, filter); err != nil {
		return
	}

	var opt = fnParams.Get(opts)
	if opt == nil {
		opt = &options.FindOptions{}
	}

	opt.Skip = fnReflect.ToPointer(p.Page() * p.Size())
	opt.Limit = fnReflect.ToPointer(p.Size())

	if ls, err = ReadAllM[DATA](ctx, i, opt); err != nil {
		return
	}

	return
}
