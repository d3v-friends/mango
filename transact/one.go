package transact

import (
	"context"
	"github.com/d3v-friends/mango/models"
	"github.com/d3v-friends/mango/mvars"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type FnTransact[MODEL any] func(model *MODEL) (update *bson.M, err error)

type SessionContext struct {
	mongo.SessionContext
}

// TransactionOne 도큐먼트 1개를 lock 해서 관리하는 트렌젝션
func TransactionOne[MODEL any](
	ctx context.Context,
	db *mongo.Database,
	iFilter models.IfFilter[MODEL],
	fn FnTransact[MODEL],
) (err error) {
	filter := iFilter.Filter()
	(*filter)[mvars.FInTrx] = false

	var single *mongo.SingleResult
	if single = db.
		Collection(iFilter.CollectionNm()).
		FindOneAndUpdate(
			ctx,
			filter,
			&bson.M{
				mvars.OSet: &bson.M{
					mvars.FInTrx: true,
				},
			},
		); single.Err() != nil {
		err = single.Err()
		return
	}

	model := iFilter.NewModel()
	if err = single.Decode(model); err != nil {
		return
	}

	var update *bson.M
	if update, err = fn(model); err != nil {
		return
	}

	(*update)[mvars.FInTrx] = true
	(*filter)[mvars.FInTrx] = true

	if single = db.Collection(iFilter.CollectionNm()).FindOneAndUpdate(ctx, filter, update); single.Err() != nil {
		err = single.Err()
		return
	}

	return
}

func createSession(db *mongo.Database) (sess mongo.Session, err error) {
	return db.Client().StartSession(&options.SessionOptions{
		DefaultReadConcern:  readconcern.Majority(),
		DefaultWriteConcern: writeconcern.Majority(),
	})
}
