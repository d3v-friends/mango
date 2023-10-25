package m_migrate

import (
	"context"
	"github.com/d3v-friends/mango/m_tx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type (
	IfMigrateModel interface {
		m_tx.IfTxModel
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
	if err = migrateDoc(ctx, db, &DocMango{}); err != nil {
		return
	}

	for _, model := range models {
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
) error {

	return m_tx.Transact(ctx, db, func(tx *m_tx.TxDB) (err error) {
		var colNm = model.GetColNm()

		var count int64
		if count, err = tx.Count(colNm, bson.M{
			"colNm": colNm,
		}); err != nil {
			return
		}

		if count == 0 {
			if err = tx.InsertOne(&DocMango{
				Id:        primitive.NewObjectID(),
				ColNm:     colNm,
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
				"colNm": colNm,
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
							Key:   m_tx.FieldInTxNm,
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

		var col = tx.Collection(colNm)
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
					"colNm": colNm,
				}, bson.M{
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
