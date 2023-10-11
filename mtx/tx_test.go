package mtx

import (
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/d3v-friends/go-pure/fnEnv"
	"github.com/d3v-friends/go-pure/fnPanic"
	"github.com/d3v-friends/mango/fn/fnMango"
	"github.com/d3v-friends/mango/mcodec"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestMTX(test *testing.T) {
	fnPanic.On(fnEnv.ReadFromFile("../env/.env"))

	var registry = mcodec.NewRegistryWithDecimal()
	var client = fnPanic.OnValue(fnMango.Connect(context.TODO(), &fnMango.IConnect{
		Host:     fnEnv.Read("MG_HOST"),
		Username: fnEnv.Read("MG_USERNAME"),
		Password: fnEnv.Read("MG_PASSWORD"),
		Registry: registry,
	}))

	var db = client.Database(fnEnv.Read("MG_DATABASE"))
	var err = fnMango.Migrate(context.TODO(), db, &testModel{})
	var tool = newTestTool(db)

	test.Run("insert one", func(t *testing.T) {
		var ctx = tool.Context()
		var tx *Transaction

		if tx, err = NewTransaction(ctx); err != nil {
			t.Fatal(err)
		}

		var model = &testModel{
			Id:        primitive.NewObjectID(),
			Name:      gofakeit.Name(),
			CreatedAt: time.Now(),
		}

		if err = tx.Transaction(func(txCtx context.Context, txDB *TxDB) (err error) {
			if err = txDB.InsertOne(ctx, model); err != nil {
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

		assert.Equal(t, 1, has)
	})

	test.Run("insert one rollback", func(t *testing.T) {

	})
}
