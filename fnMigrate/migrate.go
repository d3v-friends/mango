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
		GetMigrate() []FnMigrate
	}

	FnMigrate func(ctx context.Context, col *mongo.Collection) (memo string, err error)
)

const defMngColNm = "mango"

func Migrate(
	ctx context.Context,
	i *MigrateArgs,
) (err error) {
	var ctxKeys = make([]string, 0)
	if i.CtxKey != nil {
		ctxKeys = append(ctxKeys, *i.CtxKey)
	}

	var models = make([]Model, len(i.Models)+1)
	models[0] = &DocMango{
		colNm: i.mngColNm(),
	}
	for idx := range i.Models {
		models[idx+1] = i.Models[idx]
	}

	var db = fnMango.GetDbP(ctx, ctxKeys...)
	var colMango = db.Collection(models[0].GetColNm())
	for _, model := range models {
		var cur *mongo.SingleResult
		if cur = colMango.FindOne(ctx, bson.M{
			"colNm": model.GetColNm(),
		}); cur.Err() != nil {
			err = cur.Err()
			return
		}

		var doc *DocMango
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

		colNm string
	}

	DocMangoHistory struct {
		Memo       string    `bson:"memo"`
		MigratedAt time.Time `bson:"migratedAt"`
	}
)

func (x *DocMango) GetColNm() string {
	return x.colNm
}

func (x *DocMango) GetMigrate() []FnMigrate {
	return []FnMigrate{
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

type MigrateArgs struct {
	Models   []Model
	MngColNm *string
	CtxKey   *string
}

func (x *MigrateArgs) mngColNm() string {
	if x.MngColNm != nil {
		return *x.MngColNm
	}
	return defMngColNm
}
