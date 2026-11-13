package mgregistry_test

import (
	"testing"

	"github.com/d3v-friends/mango/v2/mgmigrate"
	"github.com/d3v-friends/mango/v2/mgquery"
	"github.com/d3v-friends/mango/v2/tester"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type DecimaModel struct {
	Id      bson.ObjectID    `bson:"_id"`
	Value   decimal.Decimal  `bson:"value"`
	Pointer *decimal.Decimal `bson:"pointer"`
}

func (x DecimaModel) GetColNm() string {
	return "decimalModels"
}

func (x DecimaModel) GetMigrates() mgmigrate.Steps {
	return mgmigrate.Steps{}
}

func TestDecimal(test *testing.T) {
	var tool = tester.NewTool(test)

	test.Run("codec", func(t *testing.T) {
		var ctx = tool.Context()
		var model = &DecimaModel{
			Id:      bson.NewObjectID(),
			Value:   decimal.Zero,
			Pointer: nil,
		}

		var _, err = mgquery.InsertOne(ctx, model)
		assert.NoError(t, err)

		model, err = mgquery.FindOne[DecimaModel](
			ctx,
			bson.M{
				"_id": model.Id,
			},
			nil,
		)
		assert.NoError(t, err)

		assert.Equal(t, "0", model.Value.String())
		assert.Nil(t, model.Pointer)

	})
}
