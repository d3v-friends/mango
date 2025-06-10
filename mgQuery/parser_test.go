package mgQuery_test

import (
	"github.com/d3v-friends/go-tools/fnPointer"
	"github.com/d3v-friends/mango/mgQuery"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"testing"
)

func TestParser(test *testing.T) {

	test.Run("sorter", func(t *testing.T) {
		var sorter, err = mgQuery.ParseSorterBsonD([]TestModelSorter{
			{
				Id: fnPointer.Make(SorterASC),
			},
			{
				Age: fnPointer.Make(SorterDESC),
			},
			{
				Data: &TestModelDataSorter{
					Title: fnPointer.Make(SorterASC),
				},
			},
		})

		assert.NoError(t, err)

		var res = bson.D{
			{
				Key: "_id", Value: int32(1),
			},
			{
				Key: "age", Value: int32(-1),
			},
			{
				Key: "data.title", Value: int32(1),
			},
		}

		assert.True(t, reflect.DeepEqual(sorter, res))
	})
}
