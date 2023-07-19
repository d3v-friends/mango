package transact

import (
	"context"
	"errors"
	"fmt"
	"github.com/d3v-friends/mango/mvars"
	"github.com/d3v-friends/pure-go/fnReflect"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"strings"
)

// One 1개의 Document 에서의 트렌젝션
// callback 함수에서 지정한 모델 이외의 데이터 조작은 최대한 지양한다
func One[MODEL any](
	ctx context.Context,
	db *mongo.Database,
	trxOne TrxOne[MODEL],
) (err error) {
	trxOne.ctx = ctx
	trxOne.col = db.Collection(trxOne.CollectionNm)

	var prevModel *MODEL
	if prevModel, err = trxOne.getModel(); err != nil {
		return
	}

	var modelId primitive.ObjectID
	if modelId, err = trxOne.getId(prevModel); err != nil {
		return
	}

	defer func() {
		if err == nil {
			return
		}
		errRollback := trxOne.rollback(modelId, prevModel)
		err = errors.Join(err, errRollback)
	}()

	var update bson.M
	if update, err = trxOne.Fn(prevModel); err != nil {
		return
	}

	err = trxOne.commit(
		modelId,
		update,
	)

	return
}

type (
	TrxOne[MODEL any] struct {
		CollectionNm string
		Filter       bson.M
		Fn           FnTrxOne[MODEL]

		// inner fields
		col *mongo.Collection
		ctx context.Context
	}

	FnTrxOne[Model any] func(model *Model) (update bson.M, err error)
)

func (x *TrxOne[MODEL]) getModel() (model *MODEL, err error) {

	x.Filter[mvars.FInTrx] = false

	var single *mongo.SingleResult
	if single = x.col.FindOneAndUpdate(
		x.ctx,
		x.Filter,
		bson.M{
			mvars.OSet: bson.M{
				mvars.FInTrx: true,
			},
		},
	); single.Err() != nil {
		err = single.Err()
		return
	}

	model = new(MODEL)
	if err = single.Decode(model); err != nil {
		return
	}

	return
}

func (x *TrxOne[MODEL]) getId(model *MODEL) (id primitive.ObjectID, err error) {
	var vo reflect.Value
	var to reflect.Type
	unbox := fnReflect.UnboxPointer(model)
	vo = reflect.ValueOf(unbox)
	to = reflect.TypeOf(unbox)
	id = primitive.NilObjectID

	for i := 0; i < vo.NumField(); i++ {
		field := vo.Field(i)
		if strings.ToLower(to.Field(i).Name) != "id" {
			continue
		}

		var isOk bool
		id, isOk = field.Interface().(primitive.ObjectID)

		if !isOk {
			continue
		}

		break
	}

	if id == primitive.NilObjectID {
		err = fmt.Errorf("not found ID from model: model=%s", vo.Type().Name())
		return
	}

	return
}

func (x *TrxOne[MODEL]) rollback(id primitive.ObjectID, prevModel *MODEL) (err error) {
	_, err = x.col.ReplaceOne(
		x.ctx,
		bson.M{
			mvars.FID:    id,
			mvars.FInTrx: true,
		},
		prevModel,
		&options.ReplaceOptions{
			Upsert: fnReflect.ToPointer(true),
		},
	)
	return
}

func (x *TrxOne[MODEL]) commit(id primitive.ObjectID, update bson.M) (err error) {
	var has bool
	if _, has = update[mvars.OSet]; has {
		switch vt := (update)[mvars.OSet].(type) {
		case bson.M:
			vt[mvars.FInTrx] = false
		case bson.D:
			vt = append(vt, bson.E{
				Key:   mvars.FInTrx,
				Value: false,
			})
		default:
			err = fmt.Errorf("invalid update value: update=%s", update)
			return
		}
	} else {
		update[mvars.OSet] = &bson.M{
			mvars.FInTrx: true,
		}
	}

	_, err = x.col.UpdateOne(
		x.ctx,
		bson.M{
			mvars.FID:    id,
			mvars.FInTrx: true,
		},
		update,
		&options.UpdateOptions{
			Upsert: fnReflect.ToPointer(true),
		},
	)

	return
}
