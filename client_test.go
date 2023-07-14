package mango

import (
	"context"
	"fmt"
	"github.com/d3v-friends/mango/migrate"
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
			"_id": primitive.NilObjectID,
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

}

type testModel struct {
	Id        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

func (x *testModel) CollectionNm() string {
	return "testModel"
}

func (x *testModel) MigrateList() migrate.FnMigrateList {
	return migrate.FnMigrateList{
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
