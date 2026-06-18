package mgmigrate_test

import (
	"context"
	"testing"

	"github.com/d3v-friends/go-tools/fnEnv"
	"github.com/d3v-friends/mango/v2/mgmigrate"
	"github.com/d3v-friends/mango/v2/tester"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const migrateMessage = "init indexing"

func TestMigrate(test *testing.T) {
	assert.NoError(test, fnEnv.Load(envPath))
	var tool = tester.NewTool(test)

	test.Run("migrate", func(t *testing.T) {
		tool.Truncate(test)
		var ctx = tool.Context()
		var err = mgmigrate.Do(
			ctx,
			tool.DB,
			MigrateTestModel{},
		)
		assert.NoError(t, err)

		var res *mongo.SingleResult
		res = tool.DB.Collection(mgmigrate.ColNm).FindOne(ctx, bson.M{
			mgmigrate.FieldColNm: MigrateTestModel{}.GetColNm(),
		})

		assert.NoError(t, res.Err())

		var migrate = new(mgmigrate.Mango)
		assert.NoError(t, res.Decode(migrate))

		assert.Equal(t, 1, migrate.NextIdx)
		assert.Equal(t, migrateMessage, migrate.History[0].Memo)

		// todo 생성된 인덱스 확인하기
	})

}

type MigrateTestModel struct {
	Id    bson.ObjectID `bson:"_id"`
	Value string        `bson:"value"`
}

func (x MigrateTestModel) GetColNm() string {
	return "migrateTestModel"
}

func (x MigrateTestModel) GetMigrates() mgmigrate.Steps {
	return mgmigrate.Steps{
		func(ctx context.Context, col *mongo.Collection) (memo string, err error) {
			memo = migrateMessage
			_, err = col.Indexes().CreateMany(ctx, []mongo.IndexModel{
				{
					Keys: bson.D{
						{Key: "value", Value: 1},
					},
				},
			})
			return
		},
	}
}
