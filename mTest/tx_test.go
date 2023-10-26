package mTest

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/d3v-friends/go-pure/fnPanic"
	"github.com/d3v-friends/mango/mTx"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

func TestTx(test *testing.T) {
	var mg = NewMango(true)

	fnPanic.On(mg.Migrate(context.TODO(), &DocTest{}))

	test.Run("insert one", func(t *testing.T) {
		var ctx = context.TODO()
		var err error
		var model = &DocTest{
			Id:        primitive.NewObjectID(),
			Name:      gofakeit.Name(),
			CreatedAt: time.Now(),
		}

		if err = mg.Tx(ctx, func(tx *mTx.TxDB) (txErr error) {
			if err = tx.InsertOne(model); err != nil {
				return
			}
			return
		}); err != nil {
			t.Fatal(err)
		}

		var has int64
		if has, err = mg.DB.
			Collection(
				model.GetColNm(),
			).
			CountDocuments(
				ctx,
				bson.M{
					"_id": model.Id,
				},
			); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, int64(1), has)
	})

	test.Run("insert one rollback", func(t *testing.T) {
		var ctx = context.TODO()
		var err error

		var model = &DocTest{
			Id:        primitive.NewObjectID(),
			Name:      gofakeit.Name(),
			CreatedAt: time.Now(),
		}

		if err = mg.Tx(ctx, func(txDB *mTx.TxDB) (err error) {
			if err = txDB.InsertOne(model); err != nil {
				return
			}
			err = fmt.Errorf("occure err to rollback")
			return
		}); err != nil {
			t.Fatal(err)
		}

		var has int64
		if has, err = mg.DB.Collection(model.GetColNm()).
			CountDocuments(
				ctx,
				bson.M{
					"_id": model.Id,
				},
			); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, int64(0), has)
	})

	test.Run("update one", func(t *testing.T) {
		var ctx = context.TODO()
		var err error

		var model = &DocTest{
			Id:        primitive.NewObjectID(),
			Name:      gofakeit.Name(),
			CreatedAt: time.Now(),
		}

		if err = mg.Tx(ctx, func(txDB *mTx.TxDB) (err error) {
			if err = txDB.InsertOne(model); err != nil {
				return
			}
			return
		}); err != nil {
			t.Fatal(err)
		}

		if err = mg.Tx(ctx, func(txDB *mTx.TxDB) (err error) {
			return txDB.UpdateOne(
				model.GetColNm(),
				bson.M{
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

		var col = mg.DB.Collection(model.GetColNm())
		var res *mongo.SingleResult
		if res = col.FindOne(
			ctx,
			bson.M{
				"_id": model.Id,
			},
		); res.Err() != nil {
			t.Fatal(res.Err())
		}

		var loadedModel = &DocTest{}
		if err = res.Decode(loadedModel); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, model.Id, loadedModel.Id)
		assert.NotEqual(t, model.Name, loadedModel.Name)
	})

	test.Run("update one rollback", func(t *testing.T) {
		var ctx = context.TODO()
		var err error
		var model = &DocTest{
			Id:        primitive.NewObjectID(),
			Name:      gofakeit.Name(),
			CreatedAt: time.Now(),
		}

		if err = mg.Tx(ctx, func(txDB *mTx.TxDB) (err error) {
			if err = txDB.InsertOne(model); err != nil {
				return
			}
			return
		}); err != nil {
			t.Fatal(err)
		}

		if err = mg.Tx(ctx, func(txDB *mTx.TxDB) (err error) {
			var name = gofakeit.BeerName()
			err = txDB.UpdateOne(
				model.GetColNm(),
				bson.M{
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
				model.GetColNm(),
				bson.M{
					"_id": model.Id,
				},
			); res.Err() != nil {
				return
			}

			var loadedModel = &DocTest{}
			if err = res.Decode(loadedModel); err != nil {
				return
			}

			assert.Equal(t, name, loadedModel.Name)

			err = fmt.Errorf("occure err for rollback")
			return
		}); err != nil {
			t.Fatal(err)
		}

		var col = mg.DB.Collection(model.GetColNm())
		var res *mongo.SingleResult
		if res = col.FindOne(ctx, bson.M{
			"_id": model.Id,
		}); res.Err() != nil {
			t.Fatal(res.Err())
		}

		var loadedModel = &DocTest{}
		if err = res.Decode(loadedModel); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, model.Id, loadedModel.Id)
		assert.Equal(t, model.Name, loadedModel.Name)
	})

	test.Run("update many", func(t *testing.T) {
		var ctx = context.TODO()
		var err error
		var try = 5

		err = mg.Tx(ctx, func(txDB *mTx.TxDB) (err error) {
			var groupId = primitive.NewObjectID()
			var now = time.Now()
			var iModels = make([]mTx.IfTxModel, try)

			for i := 0; i < try; i++ {
				iModels[i] = &DocTest{
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
			if count, err = txDB.Count(
				iModels[0].GetColNm(),
				bson.M{
					"groupId": groupId,
				},
			); err != nil {
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
		var ctx = context.TODO()
		var err error
		var try = 5

		var groupId = primitive.NewObjectID()
		var now = time.Now()
		var models = make([]mTx.IfTxModel, try)
		var iModels = make([]interface{}, try)
		var model = DocTest{}
		var colNm = model.GetColNm()
		for i := 0; i < try; i++ {
			models[i] = &DocTest{
				Id:        primitive.NewObjectID(),
				GroupId:   groupId,
				Name:      gofakeit.Name(),
				CreatedAt: now,
			}

			iModels[i] = models[i]
		}

		if _, err = mg.DB.
			Collection(colNm).
			InsertMany(ctx, iModels); err != nil {
			t.Fatal(err)
		}

		var updatedName = gofakeit.Name()
		_ = mg.Tx(ctx, func(txDB *mTx.TxDB) (err error) {
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
		if count, err = mg.DB.Collection(colNm).CountDocuments(ctx, bson.M{
			"name": updatedName,
		}); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 0, int(count))
	})

	test.Run("delete one", func(t *testing.T) {
		var ctx = context.TODO()
		var err error
		var model = &DocTest{
			Id:        primitive.NewObjectID(),
			GroupId:   primitive.NewObjectID(),
			Name:      gofakeit.Name(),
			CreatedAt: time.Now(),
		}

		if _, err = mg.DB.
			Collection(model.GetColNm()).
			InsertOne(ctx, model); err != nil {
			t.Fatal(err)
		}

		err = mg.Tx(ctx, func(txDB *mTx.TxDB) (err error) {
			return txDB.DeleteOne(model.GetColNm(), bson.M{
				"_id": model.Id,
			})
		})

		var count int64
		if count, err = mg.DB.
			Collection(model.GetColNm()).
			CountDocuments(ctx, bson.M{
				"_id": model.Id,
			}); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 0, int(count))

	})

	test.Run("delete one rollback", func(t *testing.T) {
		var ctx = context.TODO()
		var err error
		var model = &DocTest{
			Id:        primitive.NewObjectID(),
			GroupId:   primitive.NewObjectID(),
			Name:      gofakeit.Name(),
			CreatedAt: time.Now(),
		}

		if _, err = mg.DB.
			Collection(model.GetColNm()).
			InsertOne(ctx, model); err != nil {
			t.Fatal(err)
		}

		err = mg.Tx(ctx, func(txDB *mTx.TxDB) (err error) {
			err = txDB.DeleteOne(
				model.GetColNm(),
				bson.M{
					"_id": model.Id,
				},
			)

			if err != nil {
				return
			}

			err = fmt.Errorf("occure error")
			return
		})

		var count int64
		if count, err = mg.DB.
			Collection(model.GetColNm()).
			CountDocuments(ctx, bson.M{
				"_id": model.Id,
			}); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 1, int(count))

	})

	test.Run("lock test", func(t *testing.T) {
		var ctx = context.TODO()

		var id = primitive.NewObjectID()
		go mg.TxWithDelay(
			ctx,
			func(tx *mTx.TxDB) (txErr error) {
				var model = &DocTest{
					Id:        id,
					GroupId:   primitive.NewObjectID(),
					Name:      gofakeit.BeerName(),
					CreatedAt: time.Now(),
				}

				txErr = tx.InsertOne(model)

				return
			},
			5*time.Second,
		)

		var count = fnPanic.Get(mg.DB.Collection(docTestNm).CountDocuments(
			ctx,
			bson.M{
				"_id": id,
			},
		))

		assert.Equal(t, int64(1), count)

		time.Sleep(8 * time.Second)

	})
}
