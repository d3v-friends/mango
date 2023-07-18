package transact

import (
	"context"
	"github.com/d3v-friends/mango/mvars"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Transaction 1개의 Document 에서의 트렌젝션
func Transaction[MODEL any](
	ctx context.Context,
	db *mongo.Database,
	trxOne TrxOne[MODEL],
) (err error) {
	trxOne.ctx = ctx
	trxOne.col = db.Collection(trxOne.CollectionNm)

	prevModel := new(MODEL)
	if prevModel, err = trxOne.getModel(); err != nil {
		return
	}

	modelId := trxOne.getId(prevModel)

	var update *bson.M
	if update, err = trxOne.Fn(prevModel); err != nil {
		if errRollback := trxOne.rollback(
			modelId,
			prevModel,
		); errRollback != nil {
			panic(errRollback)
		}
		return
	}

	if errCommit := trxOne.commit(
		modelId,
		update,
	); errCommit != nil {
		panic(errCommit)
	}

	return
}

func getTrxUpdate(onTrx bool) (res *bson.M) {
	res = &bson.M{
		mvars.OSet: &bson.M{
			mvars.FInTrx: onTrx,
		},
	}
	return
}

type (
	TrxOne[MODEL any] struct {
		CollectionNm string
		Filter       *bson.M
		Fn           FnTrxOne[MODEL]

		// inner fields
		col *mongo.Collection
		ctx context.Context
	}

	FnTrxOne[Model any] func(model *Model) (update *bson.M, err error)
)

func (x *TrxOne[MODEL]) getModel() (res *MODEL, err error) {

	panic("not impl")
}

func (x *TrxOne[MODEL]) getId(model *MODEL) (id primitive.ObjectID) {
	panic("not impl")
}

func (x *TrxOne[MODEL]) rollback(id primitive.ObjectID, prevModel *MODEL) (err error) {
	panic("not impl")
}

func (x *TrxOne[MODEL]) commit(id primitive.ObjectID, update *bson.M) (err error) {
	panic("not impl")
}
