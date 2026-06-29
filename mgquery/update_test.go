package mgquery_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/d3v-friends/mango/v2/mgop"
	"github.com/d3v-friends/mango/v2/mgquery"
	"github.com/d3v-friends/mango/v2/tester"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestUpdate(test *testing.T) {
	var tool = tester.NewTool(test)

	test.Run("updateOne", func(t *testing.T) {
		var ctx = tool.Context()
		var sample = tool.NewSample()
		var _, err = mgquery.InsertOne(ctx, sample)
		assert.NoError(t, err)

		var name = gofakeit.Name()
		var res *mongo.UpdateResult
		res, err = mgquery.UpdateOne[tester.Sample](
			ctx,
			bson.M{
				tester.FieldId: sample.Id,
			},
			bson.M{
				mgop.Set: bson.M{
					tester.FieldName: name,
				},
			},
		)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), res.ModifiedCount)

		var loaded = new(tester.Sample)
		loaded, err = mgquery.FindOne[tester.Sample](
			ctx,
			bson.M{
				tester.FieldId: sample.Id,
			},
			nil,
		)
		assert.NoError(t, err)

		assert.Equal(t, name, loaded.Name)
	})

	test.Run("updateMany", func(t *testing.T) {
		var ctx = tool.Context()
		var samples = tool.NewSamples(10)
		var _, err = mgquery.InsertMany(ctx, samples)
		assert.NoError(t, err)

		var ids = make([]bson.ObjectID, len(samples))
		for i, sample := range samples {
			ids[i] = sample.Id
		}

		var name = gofakeit.Name()
		var res *mongo.UpdateResult
		res, err = mgquery.UpdateMany[tester.Sample](
			ctx,
			bson.M{
				tester.FieldId: bson.M{
					mgop.In: ids,
				},
			},
			bson.M{
				mgop.Set: bson.M{
					tester.FieldName: name,
				},
			},
		)

		assert.Equal(t, int64(len(samples)), res.ModifiedCount)
		var loaded = make([]*tester.Sample, len(samples))
		loaded, err = mgquery.Find[tester.Sample](
			ctx,
			bson.M{
				tester.FieldId: bson.M{
					mgop.In: ids,
				},
			},
			nil,
			nil,
		)

		for _, sample := range loaded {
			assert.Equal(t, name, sample.Name)
		}

	})

}
