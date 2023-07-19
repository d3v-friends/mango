package transact

import (
	"context"
	"fmt"
	"github.com/d3v-friends/mango"
	"github.com/d3v-friends/mango/mvars"
	"github.com/d3v-friends/pure-go/fnEnv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestOne(test *testing.T) {
	var client *mango.Client
	var err error
	if client, err = mango.NewClient(&mango.ClientOpt{
		Host:     fnEnv.Read("HOST"),
		Username: fnEnv.Read("USERNAME"),
		Password: fnEnv.Read("PASSWORD"),
		Database: fnEnv.Read("DATABASE"),
	}); err != nil {
		test.Fatal(err)
	}

	test.Run("one (not error)", func(t *testing.T) {
		ctx := context.TODO()
		now := time.Now()
		model := &testModel{
			Id:        primitive.NewObjectID(),
			Name:      "abcd",
			InTrx:     false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		var err error
		if _, err = client.
			Database().
			Collection(model.GetCollectionNm()).
			InsertOne(ctx, model); err != nil {
			t.Fatalf("%+v", err)
		}

		if err = One(
			ctx,
			client.Database(),
			TrxOne[testModel]{
				CollectionNm: model.GetCollectionNm(),
				Filter: bson.M{
					mvars.FID: model.Id,
				},
				Fn: func(model *testModel) (update bson.M, err error) {
					var count int64
					if count, err = client.Database().Collection(model.GetCollectionNm()).CountDocuments(ctx, bson.M{
						mvars.FID:    model.Id,
						mvars.FInTrx: true,
					}); err != nil {
						return
					}

					if count != 1 {
						err = fmt.Errorf("model is not locked")
						return
					}

					update = bson.M{
						mvars.OSet: bson.M{
							"name": "hello",
						},
					}
					return
				},
			}); err != nil {
			t.Fatalf("%+v", err)
		}
	})
}
