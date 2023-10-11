package migrate

import (
	"context"
	"fmt"
	"github.com/d3v-friends/mango/mtype"
	"github.com/d3v-friends/mango/mvars"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type (
	docMango struct {
		Id        primitive.ObjectID `bson:"_id"`
		InTrx     bool               `bson:"inTrx"`
		Migrate   fMigrateList       `bson:"migrate"`
		CreatedAt time.Time          `bson:"createdAt"`
		UpdatedAt time.Time          `bson:"updatedAt"`
	}

	fMigrate struct {
		Id        int       `bson:"id"`
		Name      string    `bson:"name"`
		CreatedAt time.Time `bson:"createdAt"`
	}

	fMigrateList map[string][]fMigrate
)

const (
	colMango = "mango"
)

func (x *docMango) GetID() primitive.ObjectID {
	return x.Id
}

func (x *docMango) GetCollectionNm() string {
	return colMango
}

func (x *docMango) GetMigrateList() mtype.FnMigrateList {
	return mtype.FnMigrateList{}
}

func (x *docMango) RunMigrate(
	ctx context.Context,
	db *mongo.Database,
	model mtype.IfMigrateModel,
) (err error) {

	var lsMigrateFn = model.GetMigrateList()
	var lsDone, has = x.Migrate[model.GetCollectionNm()]

	var idx = 0
	if has {
		idx = len(lsDone)
	}

	var colModel = db.Collection(model.GetCollectionNm())
	var colMigrate = db.Collection(colMango)

	for i := idx; i < len(lsMigrateFn); i++ {
		fn := lsMigrateFn[i]
		var migNm string
		if migNm, err = fn(ctx, colModel); err != nil {
			return
		}

		if _, err = colMigrate.UpdateOne(
			ctx,
			&bson.M{
				mvars.FID:    primitive.NilObjectID,
				mvars.FInTrx: true,
			},
			&bson.M{
				mvars.OPush: &bson.M{
					fmt.Sprintf("migrate.%s", model.GetCollectionNm()): &fMigrate{
						Id:        i,
						Name:      migNm,
						CreatedAt: time.Now(),
					},
				},
			},
		); err != nil {
			return
		}
	}

	return
}

func V1(
	ctx context.Context,
	db *mongo.Database,
	models ...mtype.IfModel,
) (err error) {
	var doc *docMango
	if doc, err = lockDocMangoV1(ctx, db); err != nil {
		return
	}

	defer func() {
		var unlockErr error
		if unlockErr = unlockDocMangoV1(ctx, db); unlockErr != nil {
			panic(unlockErr)
		}
	}()

	for _, model := range models {
		if err = doc.RunMigrate(ctx, db, model); err != nil {
			return
		}
	}

	return
}

func unlockDocMangoV1(
	ctx context.Context,
	db *mongo.Database,
) (err error) {
	col := db.Collection(colMango)
	if _, err = col.UpdateOne(
		ctx,
		&bson.M{
			mvars.FID:    primitive.NilObjectID,
			mvars.FInTrx: true,
		},
		&bson.M{
			mvars.OSet: &bson.M{
				mvars.FInTrx:     false,
				mvars.FUpdatedAt: time.Now(),
			},
		},
	); err != nil {
		return
	}
	return
}

func lockDocMangoV1(
	ctx context.Context,
	db *mongo.Database,
) (doc *docMango, err error) {
	var count int64
	col := db.Collection(colMango)
	now := time.Now()

	if count, err = col.CountDocuments(ctx, &bson.M{}); err != nil {
		return
	}

	if count == 0 {
		doc = &docMango{
			Id:        primitive.NilObjectID,
			InTrx:     true,
			Migrate:   make(fMigrateList),
			CreatedAt: now,
			UpdatedAt: now,
		}

		if _, err = col.InsertOne(ctx, doc); err != nil {
			return
		}

	} else {
		// todo 이곳에 isRunning 이 true 일때 어떻게 처리할지 선택할수 있게 해준다.
		var res *mongo.SingleResult
		if res = col.FindOneAndUpdate(
			ctx,
			&bson.M{
				mvars.FID:    primitive.NilObjectID,
				mvars.FInTrx: false,
			},
			&bson.M{
				mvars.OSet: &bson.M{
					mvars.FInTrx: true,
				},
			},
		); res.Err() != nil {
			err = res.Err()
			return
		}

		doc = &docMango{}
		if err = res.Decode(doc); err != nil {
			return
		}
	}

	if err = doc.RunMigrate(ctx, db, doc); err != nil {
		return
	}

	return
}
