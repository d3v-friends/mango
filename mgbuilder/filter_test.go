package mgbuilder_test

import (
	"reflect"
	"testing"

	"github.com/d3v-friends/mango/v2/mgbuilder"
	"github.com/d3v-friends/mango/v2/tester"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Filter struct {
	Id *tester.ObjectIdArgs `bson:"_id"`
}

func TestFilter(test *testing.T) {
	test.Run("filter", func(t *testing.T) {
		var filter = &Filter{
			Id: &tester.ObjectIdArgs{
				Gte: new(bson.NilObjectID),
			},
		}

		var res = mgbuilder.Filter(filter)
		assert.Equal(t, true, reflect.DeepEqual(
			res,
			bson.M{
				"_id": bson.M{
					"$gte": bson.NilObjectID,
				},
			},
		))
		t.Logf("%v\n", res)
	})
}
