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
		TxId      *primitive.ObjectID  `bson:"txId"`
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

func (x *MDoc[DATA]) Save(ctx context.Context, noUpdateTimes ...bool) (err error) {
	if x.colNm == "" {
		err = fmt.Errorf("not found col_nm in model: model=%s", fnLogger.ToJsonP(x))
		return
	}

	var noUpdateTime = fnParams.Get(noUpdateTimes)
	if !noUpdateTime {
		var now = time.Now()
		if x.origin == nil {
			x.CreatedAt = now
		}
		x.UpdatedAt = now
	}

	var col = GetColP(ctx, x.colNm)

	// 새로 작성한 모델인 경우
	if x.origin == nil {
		_, err = col.InsertOne(ctx, x)
		return
	}

	// 업데이트 모델
	var updateRes *mongo.UpdateResult
	if updateRes, err = col.
		UpdateOne(
			ctx,
			bson.M{
				"_id": x.Id,
			},
			bson.M{
				"$set": x,
			},
		); err != nil {
		return
	}

	var id, isOk = updateRes.UpsertedID.(primitive.ObjectID)
	if !isOk || id != x.Id {
		err = fmt.Errorf("fail update model: model=%s", fnLogger.ToJsonP(x))
		return
	}

	return
}

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

func UpdateValue[DATA any](origin DATA, update *DATA) DATA {
	if update != nil {
		return *update
	}
	return origin
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

func UpdateDoc[DATA any](
	ctx context.Context,
	i IfFilter,
	u *DATA,
	opts ...*options.UpdateOptions,
) (doc *MDoc[DATA], err error) {
	var filter bson.M
	if filter, err = i.Filter(); err != nil {
		return
	}

	var col = GetColP(ctx, i.ColNm())
	var now = time.Now()
	if _, err = col.UpdateOne(
		ctx,
		filter,
		bson.M{
			"$set": bson.M{
				"data":      *u,
				"updatedAt": now,
			},
			"$push": bson.M{
				"history": *u,
			},
		},
		fnParams.Get(opts),
	); err != nil {
		return
	}

	return ReadOneM[DATA](ctx, i)
}

func NewDoc[DATA any](colNm string, v *DATA) (res *MDoc[DATA]) {
	var now = time.Now()
	res = &MDoc[DATA]{
		Id:        primitive.NewObjectID(),
		Data:      v,
		History:   make([]*MDocHistory[DATA], 0),
		CreatedAt: now,
		UpdatedAt: now,
		origin:    nil,
		colNm:     colNm,
	}
	return
}
