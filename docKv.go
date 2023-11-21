package mango

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/d3v-friends/go-pure/fnReflect"
	"github.com/d3v-friends/mango/mMigrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const nmDocKv = "kvs"

type (
	DocKv[DATA any] struct {
		Id        primitive.ObjectID `bson:"_id"`
		Key       string             `bson:"key"`
		Value     []byte             `bson:"value"`
		UpdatedAt time.Time          `bson:"updatedAt"`
	}
)

func (x *DocKv[DATA]) GetID() primitive.ObjectID {
	return x.Id
}

func (x *DocKv[DATA]) GetColNm() string {
	return nmDocKv
}

func (x *DocKv[DATA]) GetMigrateList() mMigrate.FnMigrateList {
	return mgDocKv
}

func (x *DocKv[DATA]) Parse() (res *DATA, err error) {
	if len(x.Value) == 0 {
		err = fmt.Errorf("empty value")
		return
	}

	res = new(DATA)
	if err = json.Unmarshal(x.Value, res); err != nil {
		return
	}

	return
}

var mgDocKv = mMigrate.FnMigrateList{
	func(ctx context.Context, col *mongo.Collection) (memo string, err error) {
		memo = "init indexing"
		_, err = col.Indexes().CreateMany(ctx, []mongo.IndexModel{
			{
				Keys: bson.D{
					{
						Key:   "key",
						Value: 1,
					},
				},
				Options: &options.IndexOptions{
					Unique: fnReflect.ToPointer(true),
				},
			},
		})
		return
	},
}

/* ------------------------------------------------------------------------------------------------------------ */

func SetKv[DATA any](ctx context.Context, key string, value *DATA) (doc *DocKv[DATA], err error) {
	var byteValue []byte
	if byteValue, err = json.Marshal(value); err != nil {
		return
	}

	var now = time.Now()
	var col = GetColP(ctx, nmDocKv)
	var cur *mongo.SingleResult
	if cur = col.FindOneAndUpdate(
		ctx,
		bson.M{
			"key": key,
		},
		bson.M{
			"$set": bson.M{
				"value":     byteValue,
				"updatedAt": now,
			},
		},
		&options.FindOneAndUpdateOptions{
			Upsert: fnReflect.ToPointer(true),
		},
	); cur.Err() != nil {
		err = cur.Err()
		return
	}

	return GetKv[DATA](ctx, key)
}

func GetKv[DATA any](ctx context.Context, key string) (res *DocKv[DATA], err error) {
	var col = GetColP(ctx, nmDocKv)

	var singleRes *mongo.SingleResult
	if singleRes = col.FindOne(
		ctx,
		bson.M{
			"key": key,
		},
	); singleRes.Err() != nil {
		err = singleRes.Err()
		return
	}

	res = new(DocKv[DATA])
	if err = singleRes.Decode(res); err != nil {
		return
	}

	return
}
