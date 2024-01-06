package fnMigrate

import (
	"context"
	"github.com/d3v-friends/go-pure/fnReflect"
	"github.com/d3v-friends/mango/fnMango"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type (
	Model interface {
		GetColNm() string
		GetMigrate() []Run
	}

	MigrateArgs struct {
		Models []Model
	}

	Run func(ctx context.Context, col *mongo.Collection) (memo string, err error)
)

func Migrate(
	ctx context.Context,
	i *MigrateArgs,
) (err error) {
	var modelList = make([]Model, 0)
	modelList = append(modelList, &DocMango{})
	modelList = append(modelList, i.Models...)

	var db = fnMango.GetDbP(ctx)
	var colMango = db.Collection(modelList[0].GetColNm())
	var now = time.Now()
	for _, model := range modelList {
		var count int64
		if count, err = colMango.CountDocuments(ctx, bson.M{
			"colNm": model.GetColNm(),
		}); err != nil {
			return
		}

		if count == 0 {
			if _, err = colMango.InsertOne(
				ctx,
				&DocMango{
					Id:        primitive.NewObjectID(),
					ColNm:     model.GetColNm(),
					NextIdx:   0,
					History:   make([]*DocMangoHistory, 0),
					CreatedAt: now,
					UpdatedAt: now,
				},
			); err != nil {
				return
			}
		}

		var cur *mongo.SingleResult
		if cur = colMango.FindOne(
			ctx,
			bson.M{
				"colNm": model.GetColNm(),
			},
		); cur.Err() != nil {
			err = cur.Err()
			return
		}

		var doc = new(DocMango)
		if err = cur.Decode(doc); err != nil {
			return
		}

		var colModel = db.Collection(model.GetColNm())
		var migrateList = model.GetMigrate()

		for i := doc.NextIdx; i < len(migrateList); i++ {
			var fn = migrateList[i]
			var memo string
			if memo, err = fn(ctx, colModel); err != nil {
				return
			}

			if _, err = colMango.UpdateOne(
				ctx,
				bson.M{
					"colNm": model.GetColNm(),
				},
				bson.M{
					"$push": bson.M{
						"history": &DocMangoHistory{
							Memo:       memo,
							MigratedAt: time.Now(),
						},
					},
					"$inc": bson.M{
						"nextIdx": 1,
					},
				}); err != nil {
				return err
			}
		}
	}

	return
}

type (
	DocMango struct {
		Id        primitive.ObjectID `bson:"_id"`
		ColNm     string             `bson:"colNm"`
		NextIdx   int                `bson:"nextIdx"`
		History   []*DocMangoHistory `bson:"history"`
		CreatedAt time.Time          `bson:"createdAt"`
		UpdatedAt time.Time          `bson:"updatedAt"`
	}

	DocMangoHistory struct {
		Memo       string    `bson:"memo"`
		MigratedAt time.Time `bson:"migratedAt"`
	}
)

func (x *DocMango) GetColNm() string {
	return "mango"
}

func (x *DocMango) GetMigrate() []Run {
	return []Run{
		func(ctx context.Context, col *mongo.Collection) (memo string, err error) {
			memo = "init indexing"
			_, err = col.Indexes().CreateMany(ctx, []mongo.IndexModel{
				{
					Keys: bson.D{
						{
							Key:   "colNm",
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
