package mMigrate

import (
	"context"
	"github.com/d3v-friends/mango/mTx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type (
	IfMigrateModel interface {
		mTx.IfTxModel
		GetMigrateList() FnMigrateList
	}

	FnMigrate     func(ctx context.Context, col *mongo.Collection) (memo string, err error)
	FnMigrateList []FnMigrate
)

func Migrate(
	ctx context.Context,
	db *mongo.Database,
	models ...IfMigrateModel,
) (err error) {
	var ls = make([]IfMigrateModel, 0)
	ls = append(ls, &DocMango{})
	ls = append(ls, models...)

	for _, model := range ls {
		if err = migrateDoc(ctx, db, model); err != nil {
			return
		}
	}

	return
}

func migrateDoc(
	ctx context.Context,
	db *mongo.Database,
	model IfMigrateModel,
) (err error) {
	var sysCol = db.Collection(docMangoNm)

	var count int64
	if count, err = sysCol.CountDocuments(ctx, bson.M{
		"colNm": model.GetColNm(),
	}); err != nil {
		return
	}

	var now = time.Now()
	if count == 0 {
		if _, err = sysCol.InsertOne(ctx, &DocMango{
			Id:        primitive.NewObjectID(),
			ColNm:     model.GetColNm(),
			NextIdx:   0,
			Histories: make([]*DocMangoHistory, 0),
			CreatedAt: now,
		}); err != nil {
			return
		}
	}

	var singleRes *mongo.SingleResult
	if singleRes = sysCol.FindOne(ctx, bson.M{
		"colNm": model.GetColNm(),
	}); singleRes.Err() != nil {
		err = singleRes.Err()
		return
	}

	var colInfo = &DocMango{}
	if err = singleRes.Decode(colInfo); err != nil {
		return
	}

	var mgList = model.GetMigrateList()
	var modelCol = db.Collection(model.GetColNm())
	for i := colInfo.NextIdx; i < len(mgList); i++ {
		var fn = mgList[i]
		var memo string
		if memo, err = fn(ctx, modelCol); err != nil {
			return
		}

		if _, err = sysCol.UpdateOne(
			ctx,
			bson.M{
				"colNm": model.GetColNm(),
			},
			bson.M{
				"$set": bson.M{
					"nextIdx":   i + 1,
					"updatedAt": now,
				},
				"$push": bson.M{
					"histories": bson.M{
						"memo":      memo,
						"createdAt": now,
					},
				},
			}); err != nil {
			return
		}
	}

	return
}

func migrateDocByTx(
	ctx context.Context,
	db *mongo.Database,
	model IfMigrateModel,
) error {
	return mTx.Transact(ctx, db, func(tx *mTx.TxDB) (err error) {
		var migColNm = model.GetColNm()

		var count int64
		if count, err = tx.Count(
			docMangoNm,
			bson.M{
				"colNm": migColNm,
			}); err != nil {
			return
		}

		if count == 0 {
			if err = tx.InsertOne(&DocMango{
				Id:        primitive.NewObjectID(),
				ColNm:     migColNm,
				NextIdx:   0,
				Histories: make([]*DocMangoHistory, 0),
				CreatedAt: time.Now(),
			}); err != nil {
				return
			}
		}

		var loadedModal = &DocMango{}
		if err = tx.FindOneAndLock(
			docMangoNm,
			bson.M{
				"colNm": migColNm,
			},
			loadedModal,
		); err != nil {
			return
		}

		var migrateList = FnMigrateList{
			func(ctx context.Context, col *mongo.Collection) (memo string, err error) {
				memo = "inTx indexing"
				_, err = col.Indexes().CreateOne(ctx, mongo.IndexModel{
					Keys: bson.D{
						{
							Key:   mTx.FieldInTxNm,
							Value: 1,
						},
					},
				})
				return
			},
		}

		migrateList = append(migrateList, model.GetMigrateList()...)

		if len(migrateList) <= loadedModal.NextIdx {
			return
		}

		var col = tx.Collection(migColNm)
		for i := loadedModal.NextIdx; i < len(migrateList); i++ {
			var fnMigrate = migrateList[i]
			var memo string
			if memo, err = fnMigrate(ctx, col); err != nil {
				return
			}

			var now = time.Now()
			if err = tx.UpdateOneOnlyLocked(
				docMangoNm,
				bson.M{
					"colNm": migColNm,
				},
				bson.M{
					"$set": bson.M{
						"nextIdx":   i + 1,
						"updatedAt": now,
					},
					"$push": bson.M{
						"histories": bson.M{
							"memo":      memo,
							"createdAt": now,
						},
					},
				}); err != nil {
				return
			}
		}

		return
	})

}
