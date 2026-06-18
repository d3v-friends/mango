package mgbuilder_test

import (
	"reflect"
	"testing"

	"github.com/d3v-friends/mango/v2/mgbuilder"
	"github.com/d3v-friends/mango/v2/tester"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Sorter struct {
	Id   *tester.Sorter
	Name *tester.Sorter
}

func TestSorter(test *testing.T) {
	test.Run("sorter", func(t *testing.T) {
		// 순서까지 정확해야 한다.
		var sorter = []*Sorter{
			{
				Name: new(tester.SorterASC),
			},
			{
				Id: new(tester.SorterDESC),
			},
		}

		var res = mgbuilder.Sorter(sorter)
		t.Logf("%+v", res)

		assert.Equal(t, true, reflect.DeepEqual(
			res,
			bson.D{
				{
					Key:   "name",
					Value: int32(1),
				},
				{
					Key:   "_id",
					Value: int32(-1),
				},
			},
		))

	})
}
