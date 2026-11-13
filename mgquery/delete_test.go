package mgquery_test

import (
	"math/rand/v2"
	"testing"

	"github.com/d3v-friends/mango/v2/mgop"
	"github.com/d3v-friends/mango/v2/mgquery"
	"github.com/d3v-friends/mango/v2/tester"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestDelete(test *testing.T) {
	var tool = tester.NewTool(test)

	test.Run("deleteOne", func(t *testing.T) {
		var ctx = tool.Context()
		var model = tool.NewSample()
		var _, err = tool.DB.Collection(model.GetColNm()).InsertOne(ctx, model)
		assert.NoError(t, err)

		var res *mongo.DeleteResult
		res, err = mgquery.DeleteOne[tester.Sample](
			ctx,
			bson.M{
				tester.FieldId: model.Id,
			},
		)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), res.DeletedCount)

		var count int64
		count, err = mgquery.Count[tester.Sample](
			ctx,
			bson.M{
				tester.FieldId: model.Id,
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	test.Run("deleteMany", func(t *testing.T) {
		var ctx = tool.Context()
		var try = rand.IntN(10)
		var ids = make([]bson.ObjectID, try)
		for i := 0; i < try; i++ {
			var model = tool.NewSample()
			var _, err = tool.DB.Collection(model.GetColNm()).InsertOne(ctx, model)
			assert.NoError(t, err)
			ids[i] = model.Id
		}

		var res, err = mgquery.DeleteMany[tester.Sample](
			ctx,
			bson.M{
				tester.FieldId: bson.M{
					mgop.In: ids,
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, int64(try), res.DeletedCount)

	})
}
