package mdKv

import (
	"context"
	"github.com/d3v-friends/go-tools/fnPointer"
	"github.com/d3v-friends/mango/mgCtx"
	"github.com/d3v-friends/mango/mgMigrate"
	"github.com/d3v-friends/mango/mgOp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

type Model struct {
	Id        primitive.ObjectID `bson:"_id"`
	Key       string             `bson:"key"`
	Value     string             `bson:"value"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

const (
	ColNm          = "_kv"
	FieldId        = "_id"
	FieldKey       = "key"
	FieldValue     = "value"
	FieldUpdatedAt = "updatedAt"
)

var migrates = mgMigrate.Steps{
	func(ctx context.Context, col *mongo.Collection) (memo string, err error) {
		memo = "init indexing"
		_, err = col.Indexes().CreateMany(ctx, []mongo.IndexModel{
			{
				Keys: bson.D{
					{Key: FieldKey, Value: 1},
				},
				Options: &options.IndexOptions{
					Unique: fnPointer.Make(true),
				},
			},
		})
		return
	},
}

func (x Model) GetColNm() string {
	return ColNm
}

func (x Model) GetMigrates() mgMigrate.Steps {
	return migrates
}

func Find(
	ctx context.Context,
	key string,
) (model *Model, err error) {
	var col *mongo.Collection
	if col, err = mgCtx.GetCol(ctx, &Model{}); err != nil {
		return
	}

	var cur = col.FindOne(ctx, bson.M{
		FieldKey: key,
	})

	if cur.Err() != nil {
		err = cur.Err()
		return
	}

	model = &Model{}
	if err = cur.Decode(model); err != nil {
		return
	}

	return
}

func Get(
	ctx context.Context,
	key string,
) string {
	var model, err = Find(ctx, key)
	if err != nil {
		return ""
	}
	return model.Value
}

func Set(
	ctx context.Context,
	key string,
	value string,
) (err error) {
	var col *mongo.Collection
	if col, err = mgCtx.GetCol(ctx, &Model{}); err != nil {
		return
	}

	var cur = col.FindOneAndUpdate(
		ctx,
		bson.M{
			FieldKey: key,
		},
		bson.M{
			mgOp.Set: bson.M{
				FieldValue:     value,
				FieldUpdatedAt: time.Now(),
			},
		},
		&options.FindOneAndUpdateOptions{
			Upsert: fnPointer.Make(true),
		},
	)

	err = cur.Err()
	if err != nil && strings.Contains(err.Error(), "mongo: no documents in result") {
		err = nil
	}

	return
}
