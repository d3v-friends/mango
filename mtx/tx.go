package mtx

import (
	"context"
	"fmt"
	"github.com/d3v-friends/go-pure/fnParams"
	"github.com/d3v-friends/mango/mctx"
	"github.com/d3v-friends/mango/mtype"
	"github.com/d3v-friends/mango/mvars"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
)

type FnTransact func(txDB *TxDB) (err error)

const isLockColumn = "isLock"

var updateLock = bson.M{
	"$set": bson.M{
		isLockColumn: true,
	},
}

var lockFilter = bson.A{
	bson.M{
		isLockColumn: bson.M{
			mvars.OExists: false,
		},
	},
	bson.M{
		isLockColumn: false,
	},
}

var updateUnlock = bson.M{
	"$unset": bson.M{
		isLockColumn: "",
	},
}

func Transact(ctx context.Context, fn FnTransact) (err error) {
	var txDB = &TxDB{
		ctx:    ctx,
		db:     mctx.GetDBP(ctx),
		insert: make([]*insertModel, 0),
		delete: make([]*deleteModel, 0),
		update: make([]*updateModel, 0),
		lock:   make([]*lockModel, 0),
	}

	if err = fn(txDB); err == nil {
		err = txDB.commit()
	} else {
		err = txDB.rollback()
	}

	return
}

/* ------------------------------------------------------------------------------------------------------------ */

type TxDB struct {
	ctx    context.Context
	db     *mongo.Database
	insert []*insertModel
	delete []*deleteModel
	update []*updateModel
	lock   []*lockModel
}

type insertModel struct {
	colNm string
	id    primitive.ObjectID
}

type deleteModel struct {
	colNm string
	raw   bson.Raw
}

type updateModel struct {
	colNm string
	raw   bson.Raw
}

type lockModel struct {
	colNm string
	id    primitive.ObjectID
}

func (x *TxDB) commit() (err error) {
	return
}

func (x *TxDB) rollback() (err error) {

	// insert
	for _, model := range x.insert {
		var col = x.db.Collection(model.colNm)
		if _, err = col.DeleteOne(x.ctx, bson.M{
			"_id": model.id,
		}); err != nil {
			return
		}
	}

	// update
	for _, model := range x.update {
		var id = model.raw.Lookup("_id").ObjectID()
		if id == primitive.NilObjectID {
			err = fmt.Errorf("fail update rollback cuz not found _id: raw=%s", model.raw.String())
			return
		}

		var col = x.db.Collection(model.colNm)
		var filter = bson.M{
			"_id": id,
		}

		if _, err = col.UpdateOne(x.ctx, filter, bson.M{
			"$set": model.raw,
		}); err != nil {
			return
		}
	}

	// delete
	for _, model := range x.delete {
		var col = x.db.Collection(model.colNm)
		if _, err = col.InsertOne(x.ctx, model.raw); err != nil {
			return
		}
	}

	// unlock
	for _, model := range x.lock {
		var col = x.db.Collection(model.colNm)
		if _, err = col.UpdateOne(
			x.ctx,
			bson.M{
				"_id":        model.id,
				isLockColumn: false,
			}, updateUnlock,
		); err != nil {
			return
		}
	}

	return
}

func (x *TxDB) InsertOne(
	model mtype.IfModel,
) (err error) {
	var col = mctx.GetDBP(x.ctx).Collection(model.GetCollectionNm())
	if _, err = col.InsertOne(x.ctx, model); err != nil {
		return
	}

	x.insert = append(x.insert, &insertModel{
		colNm: model.GetCollectionNm(),
		id:    model.GetID(),
	})

	return
}

func (x *TxDB) InsertMany(
	models []mtype.IfModel,
) (err error) {
	if len(models) == 0 {
		err = fmt.Errorf("fail insert many, empty models")
		return
	}

	var colNm = models[0].GetCollectionNm()
	var col = x.db.Collection(colNm)

	var ls = make([]interface{}, len(models))
	for i, model := range models {
		// insert 위한 형변환
		ls[i] = model

		// rollback 을 위한 데이터 추가
		x.insert = append(x.insert, &insertModel{
			colNm: colNm,
			id:    model.GetID(),
		})
	}

	if _, err = col.InsertMany(x.ctx, ls); err != nil {
		return
	}

	return
}

func (x *TxDB) UpdateOne(
	colNm string,
	filter any,
	update any,
	opts ...*options.UpdateOptions,
) (err error) {
	var col = x.db.Collection(colNm)

	var cur *mongo.SingleResult
	if cur = col.FindOne(x.ctx, filter); cur.Err() != nil {
		err = cur.Err()
		return
	}

	var raw bson.Raw
	if raw, err = cur.DecodeBytes(); err != nil {
		return
	}

	x.update = append(x.update, &updateModel{
		colNm: colNm,
		raw:   raw,
	})

	var opt = fnParams.Get(opts)
	if _, err = col.UpdateOne(x.ctx, filter, update, opt); err != nil {
		return
	}

	return
}

// UpdateMany 1개씩 모두 불러와서 변경하므로 느려질수 있음을 인지해야 한다
func (x *TxDB) UpdateMany(
	colNm string,
	filter any,
	update any,
	opts ...*options.UpdateOptions,
) (err error) {
	var col = x.db.Collection(colNm)

	var cur *mongo.Cursor
	if cur, err = col.Find(x.ctx, filter); err != nil {
		return
	}

	for cur.Next(x.ctx) {
		var id = cur.Current.Lookup("_id").ObjectID()
		if err = x.UpdateOne(colNm, bson.M{"_id": id}, update, fnParams.Get(opts)); err != nil {
			return
		}
	}

	return
}

func (x *TxDB) DeleteOne(
	colNm string,
	filter any,
	opts ...*options.DeleteOptions,
) (err error) {
	var col = x.db.Collection(colNm)

	var res *mongo.SingleResult
	if res = col.FindOne(x.ctx, filter); res.Err() != nil {
		err = res.Err()
		return
	}

	var raw bson.Raw
	if raw, err = res.DecodeBytes(); err != nil {
		return
	}

	x.delete = append(x.delete, &deleteModel{
		colNm: colNm,
		raw:   raw,
	})

	if _, err = col.DeleteOne(x.ctx, filter, fnParams.Get(opts)); err != nil {
		return
	}

	return
}

// DeleteMany 1개씩 삭제하므로 속도가 느릴수 있다.
func (x *TxDB) DeleteMany(
	colNm string,
	filter any,
	opt *options.DeleteOptions,
) (err error) {
	var col = x.db.Collection(colNm)

	var cur *mongo.Cursor
	if cur, err = col.Find(x.ctx, filter); err != nil {
		return
	}

	var ids = make([]primitive.ObjectID, 0)
	for cur.Next(x.ctx) {
		x.delete = append(x.delete, &deleteModel{
			colNm: colNm,
			raw:   cur.Current,
		})

		// todo 이곳 체크해보기
		ids = append(ids, cur.Current.Lookup("_id").ObjectID())
	}

	// 1개씩 삭제
	for _, id := range ids {
		if err = x.DeleteOne(colNm, bson.M{"_id": id}, opt); err != nil {
			return
		}
	}

	return
}

func (x *TxDB) FindOne(
	colNm string,
	filter any,
	opts ...*options.FindOneOptions,
) (res *mongo.SingleResult) {
	return x.db.Collection(colNm).FindOne(x.ctx, filter, fnParams.Get(opts))
}

func (x *TxDB) Find(
	colNm string,
	filter any,
	opts ...*options.FindOptions,
) (res *mongo.Cursor, err error) {
	return x.db.Collection(colNm).Find(x.ctx, filter, fnParams.Get(opts))
}

func (x *TxDB) Count(
	colNm string,
	filter any,
) (count int64, err error) {
	return x.db.Collection(colNm).CountDocuments(x.ctx, filter)
}

func (x *TxDB) FindOneAndLock(
	colNm string,
	filter bson.M,
	model any,
	opts ...*options.FindOneAndUpdateOptions,
) (err error) {
	var or, has = filter[mvars.OOr]
	if has {
		var array, isOk = or.(bson.A)
		if !isOk {
			err = fmt.Errorf("filter $or is not bson.A type: $or=%s", reflect.TypeOf(or).Kind())
			return
		}

		array = append(array, lockFilter...)
	} else {
		filter[mvars.OOr] = lockFilter
	}

	var col = mctx.GetDBP(x.ctx).Collection(colNm)
	var res *mongo.SingleResult
	if res = col.FindOneAndUpdate(
		x.ctx,
		filter,
		updateLock,
		fnParams.Get(opts),
	); res.Err() != nil {
		err = res.Err()
		return
	}

	var raw bson.Raw
	if raw, err = res.DecodeBytes(); err != nil {
		return
	}

	x.lock = append(x.lock, &lockModel{
		colNm: colNm,
		id:    raw.Lookup("_id").ObjectID(),
	})

	if err = res.Decode(model); err != nil {
		return
	}

	return
}

// FindAndLock 은 MongoDB 자체에 기능이 없으므로 구현하지 않는다.
