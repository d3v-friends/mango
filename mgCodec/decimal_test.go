package mgCodec_test

import (
	"testing"

	"github.com/d3v-friends/mango/mgMigrate"
	"github.com/d3v-friends/mango/mgQuery"
	"github.com/d3v-friends/mango/tester"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DecimaModel struct {
	Id      primitive.ObjectID `bson:"_id"`
	Value   decimal.Decimal    `bson:"value"`
	Pointer *decimal.Decimal   `bson:"pointer"`
}

func (x DecimaModel) GetColNm() string {
	return "decimalModels"
}

func (x DecimaModel) GetMigrates() mgMigrate.Steps {
	return mgMigrate.Steps{}
}

func TestDecimal(test *testing.T) {
	var tool = tester.NewTool(test)

	test.Run("codec", func(t *testing.T) {
		var ctx = tool.Context()
		var model = &DecimaModel{
			Id:      primitive.NewObjectID(),
			Value:   decimal.Zero,
			Pointer: nil,
		}

		var err = mgQuery.InsertOne(ctx, model)
		assert.NoError(t, err)

		model, err = mgQuery.FindOne[DecimaModel](
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
