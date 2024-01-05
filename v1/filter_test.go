package v1

import (
	"github.com/d3v-friends/go-pure/fnReflect"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"testing"
	"time"
)

func TestTime(test *testing.T) {
	test.Run("filter", func(t *testing.T) {
		var tv = &Period{
			GT: fnReflect.ToPointer(time.Now()),
		}

		var filter, err = tv.Filter()
		assert.Equalf(t, nil, err, "err is must nil")
		assert.Equal(t, true, reflect.DeepEqual(filter, bson.M{
			"$gt": *tv.GT,
		}))

		tv = &Period{}
		filter, err = tv.Filter()
		assert.NotEqual(t, nil, err)
	})

	test.Run("ifFilter", func(t *testing.T) {

	})

}
