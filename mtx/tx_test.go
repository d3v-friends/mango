package mtx

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/d3v-friends/go-pure/fnEnv"
	"github.com/d3v-friends/go-pure/fnPanic"
	"github.com/d3v-friends/mango/fn/fnMango"
	"github.com/d3v-friends/mango/mcodec"
	"github.com/d3v-friends/mango/mctx"
	"github.com/d3v-friends/mango/mtype"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

func TestTx(test *testing.T) {
	fnPanic.On(fnEnv.ReadFromFile("../env/.env"))

	var client = fnPanic.OnValue(fnMango.Connect(context.Background(), &fnMango.IConnect{
		Host:        fnEnv.Read("MG_HOST"),
		Username:    fnEnv.Read("MG_USERNAME"),
		Password:    fnEnv.Read("MG_PASSWORD"),
		SetRegistry: mcodec.AppendDecimalCodec,
	}))

	var db = client.Database(fnEnv.Read("MG_DATABASE"))
	var err = fnMango.Migrate(context.TODO(), db, &testModel{})
	var tool = newTestTool(db)

	test.Run("insert one", func(t *testing.T) {
		var ctx = tool.Context()

		var model = &testModel{
			Id:        primitive.NewObjectID(),
			Name:      gofakeit.Name(),
			CreatedAt: time.Now(),
		}

		if err = Transact(ctx, func(txDB *TxDB) (err error) {
			if err = txDB.InsertOne(model); err != nil {
				return
			}
			return
		}); err != nil {
			t.Fatal(err)
		}

		var has int64
		if has, err = db.Collection(model.GetCollectionNm()).CountDocuments(ctx, bson.M{
			"_id": model.Id,
		}); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, int64(1), has)
	})

	test.Run("insert one rollback", func(t *testing.T) {
		var ctx = tool.Context()
		var model = &testModel{
			Id:        primitive.NewObjectID(),
			Name:      gofakeit.Name(),
			CreatedAt: time.Now(),
		}

		if err = Transact(ctx, func(txDB *TxDB) (err error) {
			if err = txDB.InsertOne(model); err != nil {
				return
			}

			err = fmt.Errorf("occure err to rollback")
			return
		}); err != nil {
			t.Fatal(err)
		}

		var has int64
		if has, err = db.Collection(model.GetCollectionNm()).CountDocuments(ctx, bson.M{
			"_id": model.Id,
		}); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, int64(0), has)
	})

	test.Run("update one", func(t *testing.T) {
		var ctx = tool.Context()
		var model = &testModel{
			Id:        primitive.NewObjectID(),
			Name:      gofakeit.Name(),
			CreatedAt: time.Now(),
		}

		if err = Transact(ctx, func(txDB *TxDB) (err error) {
			if err = txDB.InsertOne(model); err != nil {
				return
			}
			return
		}); err != nil {
			t.Fatal(err)
		}

		if err = Transact(ctx, func(txDB *TxDB) (err error) {
			return txDB.UpdateOne(
				model.GetCollectionNm(), bson.M{
					"_id": model.Id,
				},
				bson.M{
					"$set": bson.M{
						"name": gofakeit.BeerName(),
					},
				},
			)
		}); err != nil {
			t.Fatal(err)
		}

		var col = mctx.GetDBP(ctx).Collection(model.GetCollectionNm())
		var res *mongo.SingleResult
		if res = col.FindOne(ctx, bson.M{
			"_id": model.Id,
		}); res.Err() != nil {
			t.Fatal(res.Err())
		}

		var loadedModel = &testModel{}
		if err = res.Decode(loadedModel); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, model.Id, loadedModel.Id)
		assert.NotEqual(t, model.Name, loadedModel.Name)
	})

	test.Run("update one rollback", func(t *testing.T) {
		var ctx = tool.Context()
		var model = &testModel{
			Id:        primitive.NewObjectID(),
			Name:      gofakeit.Name(),
			CreatedAt: time.Now(),
		}

		if err = Transact(ctx, func(txDB *TxDB) (err error) {
			if err = txDB.InsertOne(model); err != nil {
				return
			}
			return
		}); err != nil {
			t.Fatal(err)
		}

		if err = Transact(ctx, func(txDB *TxDB) (err error) {
			var name = gofakeit.BeerName()
			err = txDB.UpdateOne(
				model.GetCollectionNm(), bson.M{
					"_id": model.Id,
				},
				bson.M{
					"$set": bson.M{
						"name": name,
					},
				},
			)

			var res *mongo.SingleResult
			if res = txDB.FindOne(
				model.GetCollectionNm(),
				bson.M{
					"_id": model.Id,
				},
			); res.Err() != nil {
				return
			}

			var loadedModel = &testModel{}
			if err = res.Decode(loadedModel); err != nil {
				return
			}

			assert.Equal(t, name, loadedModel.Name)

			err = fmt.Errorf("occure err for rollback")
			return
		}); err != nil {
			t.Fatal(err)
		}

		var col = mctx.GetDBP(ctx).Collection(model.GetCollectionNm())
		var res *mongo.SingleResult
		if res = col.FindOne(ctx, bson.M{
			"_id": model.Id,
		}); res.Err() != nil {
			t.Fatal(res.Err())
		}

		var loadedModel = &testModel{}
		if err = res.Decode(loadedModel); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, model.Id, loadedModel.Id)
		assert.Equal(t, model.Name, loadedModel.Name)
	})

	test.Run("update many", func(t *testing.T) {
		var ctx = tool.Context()
		var try = 5

		err = Transact(ctx, func(txDB *TxDB) (err error) {
			var groupId = primitive.NewObjectID()
			var now = time.Now()
			var iModels = make([]mtype.IfModel, try)
			for i := 0; i < try; i++ {
				iModels[i] = &testModel{
					Id:        primitive.NewObjectID(),
					GroupId:   groupId,
					Name:      gofakeit.Name(),
					CreatedAt: now,
				}
			}

			if err = txDB.InsertMany(iModels); err != nil {
				t.Fatal(err)
			}

			var count int64
			if count, err = txDB.Count(iModels[0].GetCollectionNm(), bson.M{
				"groupId": groupId,
			}); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, try, int(count))

			return
		})

		if err != nil {
			t.Fatal(err)
		}

	})

	test.Run("update many rollback", func(t *testing.T) {
		var ctx = tool.Context()
		var try = 5

		var groupId = primitive.NewObjectID()
		var now = time.Now()
		var models = make([]mtype.IfModel, try)
		var iModels = make([]interface{}, try)
		var model = testModel{}
		var colNm = model.GetCollectionNm()
		for i := 0; i < try; i++ {
			models[i] = &testModel{
				Id:        primitive.NewObjectID(),
				GroupId:   groupId,
				Name:      gofakeit.Name(),
				CreatedAt: now,
			}

			iModels[i] = models[i]
		}

		if _, err = tool.db.Collection(colNm).InsertMany(ctx, iModels); err != nil {
			t.Fatal(err)
		}

		var updatedName = gofakeit.Name()
		_ = Transact(ctx, func(txDB *TxDB) (err error) {
			if err = txDB.UpdateMany(
				colNm,
				bson.M{
					"groupId": groupId,
				},
				bson.M{
					"$set": bson.M{
						"name": updatedName,
					},
				},
			); err != nil {
				t.Fatal(err)
			}

			var count int64
			if count, err = txDB.Count(colNm, bson.M{
				"name": updatedName,
			}); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, try, int(count))

			err = fmt.Errorf("occur error")
			return
		})

		var count int64
		if count, err = tool.db.Collection(colNm).CountDocuments(ctx, bson.M{
			"name": updatedName,
		}); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 0, int(count))
	})

	test.Run("delete one", func(t *testing.T) {
		var ctx = tool.Context()
		var model = &testModel{
			Id:        primitive.NewObjectID(),
			GroupId:   primitive.NewObjectID(),
			Name:      gofakeit.Name(),
			CreatedAt: time.Now(),
		}

		if _, err = tool.db.
			Collection(model.GetCollectionNm()).
			InsertOne(ctx, model); err != nil {
			t.Fatal(err)
		}

		err = Transact(ctx, func(txDB *TxDB) (err error) {
			return txDB.DeleteOne(model.GetCollectionNm(), bson.M{
				"_id": model.Id,
			})
		})

		var count int64
		if count, err = tool.db.Collection(model.GetCollectionNm()).CountDocuments(ctx, bson.M{
			"_id": model.Id,
		}); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 0, int(count))

	})

	test.Run("delete one rollback", func(t *testing.T) {
		var ctx = tool.Context()
		var model = &testModel{
			Id:        primitive.NewObjectID(),
			GroupId:   primitive.NewObjectID(),
			Name:      gofakeit.Name(),
			CreatedAt: time.Now(),
		}

		if _, err = tool.db.
			Collection(model.GetCollectionNm()).
			InsertOne(ctx, model); err != nil {
			t.Fatal(err)
		}

		err = Transact(ctx, func(txDB *TxDB) (err error) {
			err = txDB.DeleteOne(model.GetCollectionNm(), bson.M{
				"_id": model.Id,
			})

			if err != nil {
				return
			}

			err = fmt.Errorf("occure error")
			return
		})

		var count int64
		if count, err = tool.db.Collection(model.GetCollectionNm()).CountDocuments(ctx, bson.M{
			"_id": model.Id,
		}); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 1, int(count))

	})
}
