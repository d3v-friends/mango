package mtrx

import (
	"context"
	"github.com/d3v-friends/mango/mctx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Transaction struct {
	db *mongo.Database
}

type FnTransaction func(mctx context.Context, txDB *TxDB) (err error)

func (x *Transaction) Transaction(ctx context.Context, fn func(mctx context.Context, txDB *TxDB) (err error)) (err error) {
	var txDB = &TxDB{
		deleteIds: make([]primitive.ObjectID, 0),
	}
	if err = fn(ctx, txDB); err != nil {
		err = txDB.Commit(ctx)
	} else {
		err = txDB.Rollback(ctx)
	}

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
	filter any
	before bson.Raw
}

type idModel struct {
	Id primitive.ObjectID `bson:"_id"`
}

type IfModel interface {
	GetID() primitive.ObjectID
}

func (x *TxDB) Commit(ctx context.Context) (err error) {
	return
}

func (x *TxDB) Rollback(ctx context.Context) (err error) {
	// delete inserted
	panic("not impl")
}

func (x *TxDB) InsertOne(
	ctx context.Context,
	colNm string,
	model IfModel,
) (err error) {
	var col = mctx.GetDBP(ctx).Collection(colNm)
	if _, err = col.InsertOne(ctx, model); err != nil {
		return
	}

	x.insert = append(x.insert, &insertModel{
		colNm: colNm,
		id:    model.GetID(),
	})

	return
}

func (x *TxDB) InsertMany(
	ctx context.Context,
	colNm string,
	models []IfModel,
) (err error) {
	var col = mctx.GetDBP(ctx).Collection(colNm)

	var ls = make([]interface{}, len(models))
	for i, model := range models {
		ls[i] = model
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
		filter: filter,
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

		// todo  두번 읽을수 있는지 확인해보기
		var i = &idModel{}
		if err = cur.Decode(i); err != nil {
			return
		}

		ids = append(ids, i.Id)
	}

	// 1개씩 삭제
	for _, id := range ids {
		if err = x.DeleteOne(ctx, colNm, bson.M{"_id": id}, opt); err != nil {
			return
		}
	}

	return
}
