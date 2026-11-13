package mgquery_test

import (
	"testing"

	"github.com/d3v-friends/mango/v2/mgop"
	"github.com/d3v-friends/mango/v2/mgquery"
	"github.com/d3v-friends/mango/v2/tester"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestInsert(test *testing.T) {
	var tool = tester.NewTool(test)

	test.Run("insertOne", func(t *testing.T) {
		var ctx = tool.Context()
		var sample = tool.NewSample()
		var _, err = mgquery.InsertOne(ctx, sample)
		assert.NoError(t, err)

		var res *mongo.SingleResult
		res = tool.DB.
			Collection(tester.ColNm).
			FindOne(ctx, bson.M{
				tester.FieldId: sample.Id,
			})
		assert.NoError(t, res.Err())

		var loaded = new(tester.Sample)
		assert.NoError(t, res.Decode(loaded))
		sample.IsSame(t, loaded)
	})

	test.Run("insertMany", func(t *testing.T) {
		var ctx = tool.Context()
		var samples = tool.NewSamples(10)
		var _, err = mgquery.InsertMany(ctx, samples)
		assert.NoError(t, err)

		var ids = make([]bson.ObjectID, len(samples))
		for i, sample := range samples {
			ids[i] = sample.Id
		}

		var res *mongo.Cursor
		res, err = tool.DB.Collection(tester.ColNm).Find(
			ctx,
			bson.M{
				tester.FieldId: bson.M{
					mgop.In: ids,
				},
			},
		)
		assert.NoError(t, err)

		var loaded = make([]*tester.Sample, len(samples))
		assert.NoError(t, res.All(ctx, &loaded))

		for i, sample := range samples {
			sample.IsSame(t, loaded[i])
		}

	})
}
