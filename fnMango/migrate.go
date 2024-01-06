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
	Models []*typ.MigrateModel
}

func Migrate(
	ctx context.Context,
	i *MigrateArgs,
) (err error) {
	var modelList = make([]*typ.MigrateModel, 0)
	modelList = append(modelList, &typ.MigrateModel{
		ColNm:   stDoc.ColNmMango,
		Migrate: stDoc.MigrateMango,
	})
	modelList = append(modelList, i.Models...)

	var db = GetDbP(ctx)
	var colMango = db.Collection(modelList[0].ColNm)
	var now = time.Now()
	for _, model := range modelList {
		var count int64
		if count, err = colMango.CountDocuments(ctx, bson.M{
			"colNm": model.ColNm,
		}); err != nil {
			return
		}

		if count == 0 {
			if _, err = colMango.InsertOne(
				ctx,
				&stDoc.Mango{
					Id:        primitive.NewObjectID(),
					ColNm:     model.ColNm,
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
				"colNm": model.ColNm,
			},
		); cur.Err() != nil {
			err = cur.Err()
			return
		}

		var doc = new(stDoc.Mango)
		if err = cur.Decode(doc); err != nil {
			return
		}

		var colModel = db.Collection(model.ColNm)
		var migrateList = model.Migrate

		for i := doc.NextIdx; i < len(migrateList); i++ {
			var fn = migrateList[i]
			var memo string
			if memo, err = fn(ctx, colModel); err != nil {
				return
			}

			if _, err = colMango.UpdateOne(
				ctx,
				bson.M{
					"colNm": model.ColNm,
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
