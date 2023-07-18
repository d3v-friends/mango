package mango

import (
	"context"
	"fmt"
	"github.com/d3v-friends/mango/models"
	"github.com/d3v-friends/mango/mvars"
	"github.com/d3v-friends/pure-go/fnEnv"
	"github.com/d3v-friends/pure-go/fnReflect"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestClient(test *testing.T) {
	client, err := NewClient(&ClientOpt{
		Host:     fnEnv.Read("HOST"),
		Username: fnEnv.Read("USERNAME"),
		Password: fnEnv.Read("PASSWORD"),
		Database: fnEnv.Read("DATABASE"),
	})

	if err != nil {
		test.Fatal(err)
	}

	test.Run("test migrate", func(t *testing.T) {
		ctx := context.TODO()
		if err = client.Migrate(ctx, &testModel{}); err != nil {
			t.Fatal(err)
		}

		var count int64
		if count, err = client.database.Collection("mango").CountDocuments(ctx, &bson.M{
			mvars.FID: primitive.NilObjectID,
			"migrate.testModel": &bson.M{
				"$exist": true,
			},
		}); err != nil {
			return
		}

		if count == 0 {
			t.Fatal(fmt.Errorf("not found migrate data"))
		}
	})

	test.Run("test transaction", func(t *testing.T) {
		//ctx := context.TODO()
		//now := time.Now()
		//
		//insertModel := &testModel{
		//	id:        primitive.NewObjectID(),
		//	Name:      fmt.Sprintf("%s", primitive.NewObjectID().Hex()),
		//	CreatedAt: now,
		//	UpdatedAt: now,
		//}
		//
		//colNm := insertModel.GetCollectionName()
		//
		//if _, err = client.Database().Collection(insertModel.GetCollectionName()).InsertOne(ctx, insertModel); err != nil {
		//	t.Fatal(err)
		//}
		//
		//if err = transact.Transaction(ctx, client.database, func(sctx *transact.SessionContext, db *mongo.Database) (fnErr error) {
		//	loadModel := &testModel{}
		//	if err = transact.FindOneAndLock(sctx, loadModel, &bson.GetFilter{
		//		mvars.FID: insertModel.id,
		//	}); err != nil {
		//		return
		//	}
		//
		//	var count int64
		//	if count, err = db.Collection(colNm).CountDocuments(ctx, &bson.GetFilter{
		//		mvars.FID:    insertModel.id,
		//		mvars.FInTrx: true,
		//	}); err != nil {
		//		return
		//	}
		//
		//	if count != 1 {
		//		err = fmt.Errorf("inTrx model is not found: inTrx=%t", true)
		//		return
		//	}
		//
		//	return
		//}); err != nil {
		//	t.Fatal(err)
		//}
		//
		//// inTrx 변경된것 확인하기
		//var count int64
		//if count, err = client.Database().Collection(colNm).CountDocuments(ctx, &bson.GetFilter{
		//	mvars.FID:    insertModel.id,
		//	mvars.FInTrx: false,
		//}); err != nil {
		//	t.Fatal(err)
		//}
		//
		//if count != 1 {
		//	err = fmt.Errorf("intrx model is not fount: inTrx=%t", false)
		//}
	})

}

type testModel struct {
	Id        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

func (x *testModel) GetCollectionNm() string {
	return "testModel"
}

func (x *testModel) GetMigrateList() models.FnMigrateList {
	return models.FnMigrateList{
		func(ctx context.Context, collection *mongo.Collection) (migrationNm string, err error) {
			migrationNm = "indexing"
			_, err = collection.Indexes().CreateOne(
				ctx,
				mongo.IndexModel{
					Keys: &bson.D{
						{
							Key:   "name",
							Value: 1,
						},
					},
					Options: &options.IndexOptions{
						Unique: fnReflect.ToPointer(true),
					},
				},
			)
			return
		},
	}
}
