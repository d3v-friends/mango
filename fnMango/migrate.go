package fnMango

import (
	"context"
	"github.com/d3v-friends/mango/stDoc"
	"github.com/d3v-friends/mango/typ"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type MigrateArgs struct {
	Models []typ.Model
}

func Migrate(
	ctx context.Context,
	i *MigrateArgs,
) (err error) {
	var modelList = make([]typ.Model, 0)
	modelList = append(modelList, &stDoc.Mango{})
	modelList = append(modelList, i.Models...)

	var db = GetDbP(ctx)
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
				&stDoc.Mango{
					Id:        primitive.NewObjectID(),
					ColNm:     model.GetColNm(),
					NextIdx:   0,
					History:   make([]*stDoc.MangoHistory, 0),
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

		var doc = new(stDoc.Mango)
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
						"history": &stDoc.MangoHistory{
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
