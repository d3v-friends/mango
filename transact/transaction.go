package transact

import (
	"fmt"
	"github.com/d3v-friends/mango/mvars"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
)

//
//import (
//	"context"
//	"fmt"
//	"github.com/d3v-friends/mango/models"
//	"github.com/d3v-friends/mango/mvars"
//	"github.com/d3v-friends/pure-go/fnParams"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/bson/primitive"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"
//	"go.mongodb.org/mongo-driver/mongo/readconcern"
//	"go.mongodb.org/mongo-driver/mongo/writeconcern"
//	"reflect"
//)
//
//type (
//	FnTransaction func(sctx *SessionContext) (fnErr error)
//)
//
//// Transact (단일 디비)가상 트렌젝션
//func Transact(
//	ctx context.Context,
//	db *mongo.Database,
//	fn FnTransaction,
//) (err error) {
//	var session mongo.Session
//	if session, err = createSession(db); err != nil {
//		return
//	}
//
//	defer session.EndSession(ctx)
//
//	panic("not impl")
//
//}
//
//// Transaction (레플리카, 샤딩 디비)
//func Transaction(
//	ctx context.Context,
//	db *mongo.Database,
//	fn FnTransaction,
//) (err error) {
//	var session mongo.Session
//	if session, err = createSession(db); err != nil {
//		return
//	}
//
//	defer session.EndSession(ctx)
//
//	if _, err = session.WithTransaction(ctx, func(sctx mongo.SessionContext) (_ any, err error) {
//		sessionContext := &SessionContext{
//			SessionContext: sctx,
//			db:             sctx.Client().Database(db.Name()),
//			trxData:        make([]*sessionFindAndLock, 0),
//		}
//
//		defer func() {
//			err = sessionContext.unlockAll()
//		}()
//
//		if err = fn(sessionContext); err != nil {
//			return
//		}
//
//		return
//	}); err != nil {
//		return
//	}
//
//	return
//}
//
//func createSession(db *mongo.Database) (sess mongo.Session, err error) {
//	return db.Client().StartSession(&options.SessionOptions{
//		DefaultReadConcern:  readconcern.Majority(),
//		DefaultWriteConcern: writeconcern.Majority(),
//	})
//}
//
//type (
//	SessionContext struct {
//		mongo.SessionContext
//		db      *mongo.Database
//		trxData []*sessionFindAndLock
//
//		// managed
//		inserted []primitive.ObjectID
//		updated  []*prevModel
//		deleted  []*prevModel
//	}
//
//	prevModel struct {
//		Id    primitive.ObjectID
//		Model any
//	}
//
//	sessionFindAndLock struct {
//		dbNm   string
//		colNm  string
//		filter *bson.Filter
//	}
//)
//
//func (x *SessionContext) Database() *mongo.Database {
//	return x.db
//}
//
//func (x *SessionContext) unlockAll() (err error) {
//	for _, data := range x.trxData {
//		col := x.Database().Collection(data.colNm)
//
//		var single *mongo.SingleResult
//		data.filter = addInTrxFilter(true, data.filter)
//		if single = col.FindOneAndUpdate(x, data.filter, &bson.Filter{
//			mvars.FInTrx: false,
//		}); single.Err() != nil {
//			err = single.Err()
//			return
//		}
//	}
//
//	return
//}
//
//func (x *SessionContext) InsertOne(model models.IfModel, iOpt ...*options.InsertOneOptions) (err error) {
//	opt := fnParams.Get(iOpt)
//	if _, err = x.db.
//		Collection(model.CollectionNm()).
//		InsertOne(x.SessionContext, model, opt); err != nil {
//		return
//	}
//
//	x.inserted = append(x.inserted, model.GetID())
//
//	return
//}
//
//func (x *SessionContext) UpdateOne(filter *bson.Filter, update any, iOpt ...*options.FindOneAndUpdateOptions) (err error) {
//	opt := fnParams.Get(iOpt)
//
//	return
//}
//
//func (x *SessionContext) FindOneAndLock() (err error) {
//
//	panic("not impl")
//}
//
//// FindOneAndLock
//// model 은 반드시 models.IfModel 인터페이스를 구현해야 한다.
//// model 은 포인터 타입이 아닌것으로 해야 한다. only struct 타입
//func FindOneAndLock(
//	sctx *SessionContext,
//	model models.IfModel,
//	filter *bson.Filter,
//	iUpdate ...*bson.Filter,
//) (err error) {
//	col := sctx.db.Collection(model.CollectionNm())
//
//	var apFilter *bson.Filter
//	apFilter = addInTrxFilter(false, filter)
//
//	var apUpdate *bson.Filter
//	if apUpdate, err = addInTrxUpdate(true, iUpdate); err != nil {
//		return
//	}
//
//	var single *mongo.SingleResult
//	if single = col.FindOneAndUpdate(
//		sctx,
//		apFilter,
//		apUpdate,
//	); single.Err() != nil {
//		err = single.Err()
//		return
//	}
//
//	if err = single.Decode(model); err != nil {
//		return
//	}
//
//	sctx.trxData = append(sctx.trxData, &sessionFindAndLock{
//		colNm:  model.CollectionNm(),
//		filter: filter,
//	})
//
//	return
//}
//

func addInTrxFilter(value bool, iFilter *bson.Filter) (res *bson.Filter) {
	(*iFilter)[mvars.FInTrx] = value
	res = iFilter
	return
}

func addInTrxUpdate(value bool, iUpdate []*bson.Filter) (res *bson.Filter, err error) {
	var update *bson.Filter
	if len(iUpdate) == 0 {
		update = &bson.Filter{}
	} else {
		update = iUpdate[0]
	}

	v, has := (*update)["$set"]
	if has {
		switch p := v.(type) {
		case *bson.Filter:
			(*p)[mvars.FInTrx] = value
		case *bson.D:
			*p = append(*p, bson.E{
				Key:   mvars.FInTrx,
				Value: value,
			})
		default:
			err = fmt.Errorf("update filter value is not supported: type=%s", reflect.TypeOf(p).Name())
			return
		}
	} else {
		(*update)["$set"] = &bson.Filter{
			mvars.FInTrx: value,
		}
	}

	res = update
	return
}
