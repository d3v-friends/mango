package mTx

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/d3v-friends/go-pure/fnPanic"
	"github.com/d3v-friends/go-pure/fnParams"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
)

type (
	TxDB struct {
		ctx    context.Context
		txId   primitive.ObjectID
		db     *mongo.Database
		insert []*insertModel
		delete []*deleteModel
		update []*updateModel
		lock   []*lockModel
	}

	IfTxModel interface {
		GetID() primitive.ObjectID
		GetColNm() string
	}

	insertModel struct {
		colNm string
		id    primitive.ObjectID
	}

	deleteModel struct {
		colNm string
		raw   bson.Raw
	}

	updateModel struct {
		colNm string
		raw   bson.Raw
	}

	lockModel struct {
		colNm string
		id    primitive.ObjectID
	}
)

const FieldInTxNm = "txId"

func NewTxDB(
	ctx context.Context,
	db *mongo.Database,
) (txDB *TxDB) {

	txDB = &TxDB{
		ctx:    ctx,
		txId:   primitive.NewObjectID(),
		db:     db,
		insert: make([]*insertModel, 0),
		delete: make([]*deleteModel, 0),
		update: make([]*updateModel, 0),
		lock:   make([]*lockModel, 0),
	}
	return
}

func (x *TxDB) unlock() (err error) {
	for _, model := range x.lock {
		var col = x.db.Collection(model.colNm)
		var res *mongo.UpdateResult
		var filter = bson.M{
			"_id":       model.id,
			FieldInTxNm: x.txId,
		}

		if res, err = col.UpdateOne(
			x.ctx,
			filter,
			bson.M{
				"$unset": bson.M{
					FieldInTxNm: "",
				},
			},
		); err != nil {
			return
		}

		if res.ModifiedCount != 1 {
			err = fmt.Errorf(
				"fail unlock document: colNm=%s, id=%s, filter=%s",
				model.colNm,
				model.id.Hex(),
				fnPanic.Get(json.Marshal(filter)),
			)
			return
		}
	}
	return
}

func (x *TxDB) commit() (err error) {
	if err = x.unlock(); err != nil {
		return
	}

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

	// update (First In Last Out) 해야 rollback 이된다.
	for i := len(x.update) - 1; i >= 0; i-- {
		var model = x.update[i]

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
	if err = x.unlock(); err != nil {
		return
	}

	return
}

func (x *TxDB) Collection(colNm string, opts ...*options.CollectionOptions) *mongo.Collection {
	return x.db.Collection(colNm, opts...)
}

func (x *TxDB) InsertOne(model IfTxModel) (err error) {
	var col = x.db.Collection(model.GetColNm())
	if _, err = col.InsertOne(x.ctx, model); err != nil {
		return
	}

	x.insert = append(x.insert, &insertModel{
		colNm: model.GetColNm(),
		id:    model.GetID(),
	})

	return
}

func (x *TxDB) InsertMany(models []IfTxModel) (err error) {
	if len(models) == 0 {
		err = fmt.Errorf("fail insert many, empty models")
		return
	}

	var colNm = models[0].GetColNm()
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

func (x *TxDB) UpdateOneOnlyLocked(
	colNm string,
	filter bson.M,
	update bson.M,
	opts ...*options.UpdateOptions,
) (err error) {
	filter[FieldInTxNm] = x.txId
	return x.UpdateOne(colNm, filter, update, opts...)
}

func (x *TxDB) UpdateOne(
	colNm string,
	filter bson.M,
	update bson.M,
	opts ...*options.UpdateOptions,
) (err error) {
	var col = x.db.Collection(colNm)

	var cur *mongo.SingleResult
	if cur = col.FindOne(x.ctx, filter); cur.Err() != nil {
		err = cur.Err()
		return
	}

	var raw bson.Raw
	if raw, err = cur.Raw(); err != nil {
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

func (x *TxDB) UpdateManyOnlyLocked(
	colNm string,
	filter bson.M,
	update bson.M,
	opts ...*options.UpdateOptions,
) (err error) {
	filter[FieldInTxNm] = x.txId
	return x.UpdateMany(colNm, filter, update, opts...)
}

// UpdateMany 1개씩 모두 불러와서 변경하므로 느려질수 있음을 인지해야 한다
func (x *TxDB) UpdateMany(
	colNm string,
	filter bson.M,
	update bson.M,
	opts ...*options.UpdateOptions,
) (err error) {
	var col = x.db.Collection(colNm)

	var cur *mongo.Cursor
	if cur, err = col.Find(x.ctx, filter); err != nil {
		return
	}

	for cur.Next(x.ctx) {
		var id = cur.Current.Lookup("_id").ObjectID()
		if err = x.UpdateOne(
			colNm,
			bson.M{"_id": id},
			update,
			fnParams.Get(opts),
		); err != nil {
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
	model IfTxModel,
	opts ...*options.FindOneAndUpdateOptions,
) (err error) {
	var orFilter = bson.A{
		bson.M{
			FieldInTxNm: bson.M{
				"$exists": false,
			},
		},
	}

	var or, has = filter["$or"]
	if has {
		var array, isOk = or.(bson.A)
		if !isOk {
			err = fmt.Errorf("filter $or is not bson.A type: $or=%s", reflect.TypeOf(or).Kind())
			return
		}

		array = append(array, orFilter...)
	} else {
		filter["$or"] = orFilter
	}

	var col = x.db.Collection(colNm)
	var res *mongo.SingleResult
	if res = col.FindOneAndUpdate(
		x.ctx,
		filter,
		bson.M{
			"$set": bson.M{
				FieldInTxNm: x.txId,
			},
		},
		fnParams.Get(opts),
	); res.Err() != nil {
		err = res.Err()
		return
	}

	var raw bson.Raw
	if raw, err = res.Raw(); err != nil {
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
