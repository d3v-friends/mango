package cores

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type (
	IfInitializeModel interface {
		CollectionNm() string
		InitializeList() []FnInitializeJob
	}

	FnInitializeJob func(ctx context.Context, db *mongo.Collection) (jobNm string, err error)
)

const (
	colMangoSystem = "mangoSystem"
)

func Initialize(
	ctx context.Context,
	db *mongo.Database,
	modelList ...IfInitializeModel,
) (err error) {
	var system *DocMangoSystem
	if system, err = readMongoSystem(ctx, db); err != nil {
		return
	}

	for _, model := range modelList {
		var idx = 0
		var has bool
		var doneList FInitList
		var colNm = model.CollectionNm()

		if doneList, has = system.InitMap[colNm]; has {
			idx = len(doneList) - 1
		}

		var initList []FnInitializeJob
		initList = model.InitializeList()
		if len(initList) == 0 {
			continue
		}

		var col *mongo.Collection
		col = db.Collection(colNm)
		for i := idx; i < len(initList); i++ {
			var jobNm string
			if jobNm, err = initList[i](ctx, col); err != nil {
				return
			}

			if err = pushInitJobMongoSystem(
				ctx,
				db,
				colNm,
				&FInit{
					Idx:   i,
					JobNm: jobNm,
				}); err != nil {
				return
			}
		}
	}

	return
}

type (
	DocMangoSystem struct {
		Id        primitive.ObjectID `bson:"_id"`
		InitMap   FInitMap           `bson:"initMap"`
		CreatedAt time.Time          `bson:"createdAt"`
		UpdatedAt time.Time          `bson:"updatedAt"`
	}

	FInit struct {
		Idx   int    `bson:"idx"`
		JobNm string `bson:"jobNm"`
	}

	FInitList []FInit

	FInRunning struct {
		Init bool `bson:"init"`
	}

	FInitMap map[string][]FInit
)

func readMongoSystem(
	ctx context.Context,
	db *mongo.Database,
) (res *DocMangoSystem, err error) {
	now := time.Now()
	var cur *mongo.SingleResult
	cur = db.Collection(colMangoSystem).FindOne(ctx, &bson.M{
		"_id": primitive.NilObjectID,
	})

	switch cur.Err() {
	case mongo.ErrNoDocuments:
		res = &DocMangoSystem{
			Id:        primitive.NilObjectID,
			InitMap:   make(FInitMap),
			CreatedAt: now,
			UpdatedAt: now,
		}

		_, err = db.Collection(colMangoSystem).InsertOne(ctx, res)
		return
	case nil:
		res = &DocMangoSystem{}
		err = cur.Decode(res)
		return
	default:
		return
	}
}

func pushInitJobMongoSystem(
	ctx context.Context,
	db *mongo.Database,
	collectionNm string,
	data *FInit,
) (err error) {
	_, err = db.Collection(colMangoSystem).UpdateOne(ctx,
		&bson.M{
			"_id": primitive.NilObjectID,
		},
		&bson.M{
			"$push": &bson.M{
				fmt.Sprintf("initJob.%s", collectionNm): data,
			},
		},
	)
	return
}
