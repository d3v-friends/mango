package mgMigrate

import (
	"context"
	"github.com/d3v-friends/go-tools/fnPointer"
	"github.com/d3v-friends/go-tools/fnSlice"
	"github.com/d3v-friends/mango"
	"github.com/d3v-friends/mango/mgCtx"
	"github.com/d3v-friends/mango/mgOp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type (
	Mango struct {
		Id        primitive.ObjectID `bson:"_id"`
		ColNm     string             `bson:"colNm"`
		NextIdx   int                `bson:"nextIdx"`
		History   []*MangoHistory    `bson:"history"`
		CreatedAt time.Time          `bson:"createdAt"`
		UpdatedAt time.Time          `bson:"updatedAt"`
	}

	MangoHistory struct {
		Memo       string    `bson:"memo"`
		MigratedAt time.Time `bson:"migratedAt"`
	}

	MigratedModel interface {
		mango.Model
		GetMigrates() Steps
	}

	Step func(ctx context.Context, col *mongo.Collection) (memo string, err error)

	Steps []Step
)

const (
	FieldColNm   = "colNm"
	FieldHistory = "history"
	FieldNextIdx = "nextIdx"
	ColNm        = "_mango"
)

var migrates = Steps{
	func(ctx context.Context, collection *mongo.Collection) (memo string, err error) {
		memo = "init indexing"
		_, err = collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
			{
				Keys: bson.D{
					{Key: FieldColNm, Value: 1},
				},
				Options: &options.IndexOptions{
					Unique: fnPointer.Make(true),
				},
			},
		})
		return
	},
}

func (x *Mango) GetColNm() string {
	return ColNm
}

func (x *Mango) GetMigrates() Steps {
	return migrates
}

/*------------------------------------------------------------------------------------------------*/

func Do(
	ctx context.Context,
	db *mongo.Database,
	models ...MigratedModel,
) (err error) {
	ctx = mgCtx.SetDB(ctx, db)

	var mangoModel = &Mango{}
	models = reverseMigrateModels(append(models, mangoModel))
	if err = createAllCollection(ctx, db, models); err != nil {
		return
	}

	var now = time.Now()
	var colMango = db.Collection(mangoModel.GetColNm())
	for _, model := range models {
		var count int64
		if count, err = colMango.CountDocuments(ctx, bson.M{
			FieldColNm: model.GetColNm(),
		}); err != nil {
			return
		}

		if count == 0 {
			if _, err = colMango.InsertOne(
				ctx,
				&Mango{
					Id:        primitive.NewObjectID(),
					ColNm:     model.GetColNm(),
					NextIdx:   0,
					History:   make([]*MangoHistory, 0),
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
				FieldColNm: model.GetColNm(),
			},
		); cur.Err() != nil {
			err = cur.Err()
			return
		}

		var doc = new(Mango)
		if err = cur.Decode(doc); err != nil {
			return
		}

		var colModel = db.Collection(model.GetColNm())
		var migrateList = model.GetMigrates()

		for i := doc.NextIdx; i < len(migrateList); i++ {
			var fn = migrateList[i]
			var memo string
			if memo, err = fn(ctx, colModel); err != nil {
				return
			}

			if _, err = colMango.UpdateOne(
				ctx,
				bson.M{
					FieldColNm: model.GetColNm(),
				},
				bson.M{
					mgOp.Push: bson.M{
						FieldHistory: &MangoHistory{
							Memo:       memo,
							MigratedAt: time.Now(),
						},
					},
					mgOp.Inc: bson.M{
						FieldNextIdx: 1,
					},
				}); err != nil {
				return err
			}
		}
	}

	return
}

func createAllCollection(
	ctx context.Context,
	db *mongo.Database,
	models []MigratedModel,
) (err error) {
	var names []string
	if names, err = db.ListCollectionNames(ctx, bson.M{}); err != nil {
		return
	}

	for _, model := range models {
		var has bool
		if has = fnSlice.Has(names, func(v string) bool {
			return v == model.GetColNm()
		}); has {
			continue
		}

		if err = db.CreateCollection(ctx, model.GetColNm()); err != nil {
			return
		}
	}

	return
}

func reverseMigrateModels(ls []MigratedModel) (res []MigratedModel) {
	res = make([]MigratedModel, len(ls))
	for i := 0; i < len(ls); i++ {
		var idx = len(ls) - (i + 1)
		res[idx] = ls[i]
	}
	return
}
