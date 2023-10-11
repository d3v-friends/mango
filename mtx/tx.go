package mtx

import (
	"context"
	"fmt"
	"github.com/d3v-friends/mango/mctx"
	"github.com/d3v-friends/mango/mtype"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Transaction struct {
	ctx context.Context
}

type FnTransaction func(txCtx context.Context, txDB *TxDB) (err error)

func NewTransaction(ctx context.Context) (res *Transaction, err error) {
	res = &Transaction{
		ctx: ctx,
	}
	return
}

func (x *Transaction) Transaction(fn func(txCtx context.Context, txDB *TxDB) (err error)) (err error) {
	var txDB = &TxDB{
		insert: make([]*insertModel, 0),
		delete: make([]*deleteModel, 0),
		update: make([]*updateModel, 0),
	}

	if err = fn(x.ctx, txDB); err != nil {
		err = txDB.Commit(x.ctx)
	} else {
		err = txDB.Rollback(x.ctx)
	}

	return
}

/* ------------------------------------------------------------------------------------------------------------ */

type TxDB struct {
	insert []*insertModel
	delete []*deleteModel
	update []*updateModel
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
	colNm  string
	before bson.Raw
}

type idModel struct {
	Id primitive.ObjectID `bson:"_id"`
}

func (x *TxDB) Commit(ctx context.Context) (err error) {
	return
}

func (x *TxDB) Rollback(ctx context.Context) (err error) {
	var db = mctx.GetDBP(ctx)

	// insert
	for _, model := range x.insert {
		var col = db.Collection(model.colNm)
		if _, err = col.DeleteOne(ctx, bson.M{
			"_id": model.id,
		}); err != nil {
			return
		}
	}

	// update
	for _, model := range x.update {
		var id = model.before.Lookup("_id").ObjectID()
		if id == primitive.NilObjectID {
			err = fmt.Errorf("fail update rollback cuz not found _id: raw=%s", model.before.String())
			return
		}

		var col = db.Collection(model.colNm)
		var filter = bson.M{
			"_id": id,
		}

		if _, err = col.UpdateOne(ctx, filter, model.before); err != nil {
			return
		}
	}

	// delete
	for _, model := range x.delete {
		var col = db.Collection(model.colNm)
		if _, err = col.InsertOne(ctx, model.raw); err != nil {
			return
		}
	}

	return
}

func (x *TxDB) InsertOne(
	ctx context.Context,
	model mtype.IfModel,
) (err error) {
	var col = mctx.GetDBP(ctx).Collection(model.GetCollectionNm())
	if _, err = col.InsertOne(ctx, model); err != nil {
		return
	}

	x.insert = append(x.insert, &insertModel{
		colNm: model.GetCollectionNm(),
		id:    model.GetID(),
	})

	return
}

func (x *TxDB) InsertMany(
	ctx context.Context,
	models []mtype.IfModel,
) (err error) {
	if len(models) == 0 {
		err = fmt.Errorf("fail insert many, empty models")
		return
	}

	var db = mctx.GetDBP(ctx)
	var colNm = models[0].GetCollectionNm()
	var col = db.Collection(colNm)

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

	if _, err = col.InsertMany(ctx, ls); err != nil {
		return
	}

	return
}

func (x *TxDB) UpdateOne(
	ctx context.Context,
	colNm string,
	filter any,
	update any,
	opt *options.UpdateOptions,
) (err error) {
	var col = mctx.GetDBP(ctx).Collection(colNm)

	var cur *mongo.SingleResult
	if cur = col.FindOne(ctx, filter); cur.Err() != nil {
		err = cur.Err()
		return
	}

	var raw bson.Raw
	if raw, err = cur.DecodeBytes(); err != nil {
		return
	}

	x.update = append(x.update, &updateModel{
		colNm:  colNm,
		before: raw,
	})

	if _, err = col.UpdateOne(ctx, filter, update, opt); err != nil {
		return
	}

	return
}

// UpdateMany 1개씩 모두 불러와서 변경하므로 느려질수 있음을 인지해야 한다
func (x *TxDB) UpdateMany(
	ctx context.Context,
	colNm string,
	filter any,
	update any,
	opt *options.UpdateOptions,
) (err error) {

	var col = mctx.GetDBP(ctx).Collection(colNm)

	var cur *mongo.Cursor
	if cur, err = col.Find(ctx, filter); err != nil {
		return
	}

	var ids = make([]*idModel, 0)
	if err = cur.All(ctx, &ids); err != nil {
		return
	}

	// 1개씩 업데이트
	for _, id := range ids {
		if err = x.UpdateOne(ctx, colNm, bson.M{"_id": id}, update, opt); err != nil {
			return
		}
	}

	return
}

func (x *TxDB) DeleteOne(
	ctx context.Context,
	colNm string,
	filter any,
	opt *options.DeleteOptions,
) (err error) {
	var col = mctx.GetDBP(ctx).Collection(colNm)

	var res *mongo.SingleResult
	if res = col.FindOne(ctx, filter); res.Err() != nil {
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

	if _, err = col.DeleteOne(ctx, filter, opt); err != nil {
		return
	}

	return
}

// DeleteMany 1개씩 삭제하므로 속도가 느릴수 있다.
func (x *TxDB) DeleteMany(
	ctx context.Context,
	colNm string,
	filter any,
	opt *options.DeleteOptions,
) (err error) {
	var col = mctx.GetDBP(ctx).Collection(colNm)

	var cur *mongo.Cursor
	if cur, err = col.Find(ctx, filter); err != nil {
		return
	}

	var ids = make([]primitive.ObjectID, 0)
	for cur.Next(ctx) {
		x.delete = append(x.delete, &deleteModel{
			colNm: colNm,
			raw:   cur.Current,
		})

		// todo 이곳 체크해보기
		ids = append(ids, cur.Current.Lookup("_id").ObjectID())
	}

	// 1개씩 삭제
	for _, id := range ids {
		if err = x.DeleteOne(ctx, colNm, bson.M{"_id": id}, opt); err != nil {
			return
		}
	}

	return
}
