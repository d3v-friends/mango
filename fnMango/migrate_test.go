package fnMango

import (
	"context"
	"github.com/d3v-friends/go-pure/fnEnv"
	"github.com/d3v-friends/go-pure/fnPanic"
	"github.com/d3v-friends/go-pure/fnReflect"
	"github.com/d3v-friends/mango/stDoc"
	"github.com/d3v-friends/mango/typ"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestMigrate(test *testing.T) {
	fnPanic.On(fnEnv.ReadFromFile("../env/.env"))
	var client = fnPanic.OnValue(Connect(&ConnectArgs{
		Host:     fnEnv.Read("MG_HOST"),
		Username: fnEnv.Read("MG_USERNAME"),
		Password: fnEnv.Read("MG_PASSWORD"),
		SetRegistry: []FnSetRegistry{
			DecimalRegistry,
		},
	}))

	test.Run("migrate", func(t *testing.T) {
		var ctx = context.TODO()
		ctx = SetDb(ctx, client.Database(fnEnv.Read("MG_DATABASE")))
		var err = Migrate(ctx, &MigrateArgs{
			Models: []*typ.MigrateModel{
				{
					ColNm:   stDoc.ColNmMango,
					Migrate: stDoc.MigrateMango,
				},
			},
		})

		if err != nil {
			t.Fatal(err)
		}
	})
}

type DocTest struct {
	Id        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

func (x *DocTest) GetColNm() string {
	return "tests"
}

func (x *DocTest) GetMigrate() []typ.FnMigrate {
	return []typ.FnMigrate{
		func(ctx context.Context, col *mongo.Collection) (memo string, err error) {
			memo = "init indexing"
			_, err = col.Indexes().CreateMany(ctx, []mongo.IndexModel{
				{
					Keys: bson.D{
						{
							Key:   "name",
							Value: 1,
						},
					},
					Options: &options.IndexOptions{
						Unique: fnReflect.ToPointer(true),
					},
				},
			})
			return
		},
	}
}
