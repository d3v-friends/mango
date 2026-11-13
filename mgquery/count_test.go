package mgquery_test

import (
	"math/rand/v2"
	"testing"

	"github.com/d3v-friends/mango/v2/mgquery"
	"github.com/d3v-friends/mango/v2/tester"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestCount(test *testing.T) {
	var tool = tester.NewTool(test)
	tool.Truncate(test)
	tool.Migrate(test)

	test.Run("count", func(test *testing.T) {
		var try = rand.IntN(10)
		var ctx = tool.Context()

		for i := 0; i < try; i++ {
			var model = tool.NewSample()
			var _, err = tool.DB.Collection(model.GetColNm()).InsertOne(ctx, model)
			assert.NoError(test, err)
		}

		var count, err = mgquery.Count[tester.Sample](ctx, bson.M{})
		assert.NoError(test, err)
		assert.Equal(test, try, int(count))

	})
}
