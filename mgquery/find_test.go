package mgquery_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/d3v-friends/go-tools/fnSlice"
	"github.com/d3v-friends/mango/v2/mgop"
	"github.com/d3v-friends/mango/v2/mgquery"
	"github.com/d3v-friends/mango/v2/tester"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestFind(test *testing.T) {
	var tool = tester.NewTool(test)
	tool.Truncate(test)
	tool.Migrate(test)

	test.Run("findOne", func(t *testing.T) {
		var ctx = tool.Context()
		var try = 10
		var models = make([]*tester.Sample, try)
		for i := 0; i < try; i++ {
			models[i] = tool.NewSample()
			var _, err = mgquery.InsertOne(ctx, models[i])
			assert.NoError(t, err)
		}

		for _, model := range models {
			var loaded, err = mgquery.FindOne[tester.Sample](
				ctx,
				bson.M{
					tester.FieldId: model.Id,
				},
				nil,
			)

			assert.NoError(t, err)
			model.IsSame(test, loaded)
		}
	})

	test.Run("find", func(t *testing.T) {
		var ctx = tool.Context()
		var capacity = 50
		var models = tool.NewSamples(capacity)
		var _, err = mgquery.InsertMany(ctx, models)
		assert.NoError(t, err)

		var picked = fnSlice.PickRand(models, 20)
		var ids = fnSlice.Each[*tester.Sample, bson.ObjectID](picked, func(v *tester.Sample) bson.ObjectID {
			return v.Id
		})

		var loadedSamples []*tester.Sample
		loadedSamples, err = mgquery.Find[tester.Sample](
			ctx,
			bson.M{
				tester.FieldId: bson.M{
					mgop.In: ids,
				},
			},
			nil,
			nil,
		)

		for _, model := range picked {
			var sample *tester.Sample
			sample, err = fnSlice.FindOne(loadedSamples, func(v *tester.Sample) bool {
				return v.Id.Hex() == model.Id.Hex()
			})
			assert.NoError(t, err)

			model.IsSame(test, sample)
		}
	})

	test.Run("findList", func(t *testing.T) {
		tool.TruncateSampleCollection(t)

		var ctx = tool.Context()
		var capacity = 50
		var models = tool.NewSamples(capacity)
		var _, err = mgquery.InsertMany(ctx, models)
		assert.NoError(t, err)

		for i := 0; i < 5; i++ {
			var list *mgquery.List[tester.Sample]
			list, err = mgquery.FindList[tester.Sample](
				ctx,
				bson.M{
					tester.FieldId: bson.M{
						mgop.Exists: true,
					},
				},
				bson.D{
					{
						Key:   tester.FieldId,
						Value: 1,
					},
				},
				tester.Pager{
					Page: int64(i),
					Size: 10,
				},
			)
			assert.NoError(t, err)

			for i2, model := range list.List {
				var origin = models[i*10+i2]
				model.IsSame(t, origin)
			}
		}
	})

	test.Run("findOneAndUpdate", func(t *testing.T) {
		var ctx = tool.Context()

		var model = tool.NewSample()
		var _, err = mgquery.InsertOne(ctx, model)
		assert.NoError(t, err)

		var name = gofakeit.Name()
		var updatedModel *tester.Sample
		updatedModel, err = mgquery.FindOneAndUpdate[tester.Sample](
			ctx,
			bson.M{
				tester.FieldId: model.Id,
			},
			nil,
			bson.M{
				mgop.Set: bson.M{
					tester.FieldName: name,
				},
			},
		)
		assert.NoError(t, err)

		assert.Equal(t, model.Id.Hex(), updatedModel.Id.Hex())
		assert.Equal(t, name, updatedModel.Name)

	})
}
