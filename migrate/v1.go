package migrate

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type (
	docMango struct {
		Id        primitive.ObjectID `bson:"_id"`
		IsRunning bool               `bson:"isRunning"`
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

	IfMigrateModel interface {
		CollectionNm() string
		MigrateList() FnMigrateList
	}

	FnMigrate     func(ctx context.Context, collection *mongo.Collection) (migrationNm string, err error)
	FnMigrateList []FnMigrate
)

const (
	colMango = "mango"
)

func (x *docMango) CollectionNm() string {
	return colMango
}

func (x *docMango) MigrateList() FnMigrateList {
	return FnMigrateList{}
}

func (x *docMango) RunMigrate(
	ctx context.Context,
	db *mongo.Database,
	model IfMigrateModel,
) (err error) {

	lsMigrateFn := model.MigrateList()

	lsDone, has := x.Migrate[model.CollectionNm()]
	idx := 0
	if has {
		idx = len(lsDone) - 1
	}

	colModel := db.Collection(model.CollectionNm())
	colMigrate := db.Collection(colMango)

	for i := idx; i < len(lsMigrateFn); i++ {
		fn := lsMigrateFn[i]
		var migNm string
		if migNm, err = fn(ctx, colModel); err != nil {
			return
		}

		if _, err = colMigrate.UpdateOne(
			ctx,
			&bson.M{
				"_id":       primitive.NilObjectID,
				"isRunning": true,
			},
			&bson.M{
				"$push": &bson.M{
					fmt.Sprintf("migrate.%s", model.CollectionNm()): &fMigrate{
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
	models ...IfMigrateModel,
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
			"_id":       primitive.NilObjectID,
			"isRunning": true,
		},
		&bson.M{
			"$set": &bson.M{
				"isRunning": false,
				"updatedAt": time.Now(),
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
			IsRunning: true,
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
				"_id":       primitive.NilObjectID,
				"isRunning": false,
			},
			&bson.M{
				"$set": &bson.M{
					"isRunning": true,
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
