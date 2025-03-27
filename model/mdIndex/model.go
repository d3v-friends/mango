package mdIndex

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
	"time"
)

type Model struct {
	Id        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Number    uint64             `bson:"number"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

const (
	ColNm          = "_indexes"
	FieldId        = "_id"
	FieldName      = "name"
	FieldNumber    = "number"
	FieldUpdatedAt = "updatedAt"
)

var migrates = mgMigrate.Steps{
	func(ctx context.Context, col *mongo.Collection) (memo string, err error) {
		memo = "init indexing"
		_, err = col.Indexes().CreateMany(ctx, []mongo.IndexModel{
			{
				Keys: bson.D{
					{Key: FieldName, Value: 1},
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

func NewIndexNumber(
	ctx context.Context,
	name string,
) (idx uint64, err error) {
	var col *mongo.Collection
	if col, err = mgCtx.GetCol(ctx, &Model{}); err != nil {
		return
	}

	var res = col.FindOneAndUpdate(
		ctx,
		bson.M{
			FieldName: name,
		},
		bson.M{
			mgOp.Set: bson.M{
				FieldUpdatedAt: time.Now(),
			},
			mgOp.Inc: bson.M{
				FieldNumber: 1,
			},
		},
		&options.FindOneAndUpdateOptions{
			ReturnDocument: fnPointer.Make(options.After),
			Upsert:         fnPointer.Make(true),
		},
	)

	if res.Err() != nil {
		err = res.Err()
		return
	}

	var model = &Model{}
	if err = res.Decode(model); err != nil {
		return
	}

	idx = model.Number
	return
}
