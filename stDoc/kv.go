package stDoc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/d3v-friends/go-pure/fnReflect"
	"github.com/d3v-friends/mango/fnMango"
	"github.com/d3v-friends/mango/typ"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Kv struct {
	Id        primitive.ObjectID `bson:"_id"`
	Key       string             `bson:"key"`
	Value     []byte             `bson:"value"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

func (x *Kv) GetColNm() string {
	return "kvs"
}

func (x *Kv) GetMigrate() []typ.FnMigrate {
	return []typ.FnMigrate{
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
}

func GetKv[T any](ctx context.Context, key string, defs ...*T) (res *T, err error) {
	var now = time.Now()
	var doc = new(Kv)
	var col = fnMango.GetDbP(ctx, key).Collection(doc.GetColNm())
	var total int64
	if total, err = col.CountDocuments(
		ctx,
		bson.M{
			"key": key,
		},
	); err != nil {
		return
	}

	if total == 0 {
		if len(defs) == 0 {
			err = fmt.Errorf("not found kv")
		}

		var value []byte
		if value, err = json.Marshal(defs[0]); err != nil {
			return
		}

		if _, err = col.InsertOne(ctx, &Kv{
			Id:        primitive.NewObjectID(),
			Key:       key,
			Value:     value,
			UpdatedAt: now,
		}); err != nil {
			return
		}

		res = defs[0]
		return
	}

	var cur *mongo.SingleResult
	if cur = col.FindOne(
		ctx,
		bson.M{
			"key": key,
		},
	); cur.Err() != nil {
		err = cur.Err()
		return
	}

	if err = cur.Decode(doc); err != nil {
		return
	}

	res = new(T)
	if err = json.Unmarshal(doc.Value, res); err != nil {
		return
	}

	return
}

func SetKv[T any](ctx context.Context, key string, value *T) (err error) {
	var doc = new(Kv)
	var now = time.Now()
	var byteValue []byte
	if byteValue, err = json.Marshal(value); err != nil {
		return
	}
	if _, err = fnMango.GetDbP(ctx).Collection(doc.GetColNm()).UpdateOne(
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
		&options.UpdateOptions{
			Upsert: fnReflect.ToPointer(true),
		}); err != nil {
		return
	}

	return
}
