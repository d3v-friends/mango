package mango

import (
	"fmt"
	"github.com/d3v-friends/go-pure/fnPanic"
	"github.com/d3v-friends/go-pure/fnParams"
	"github.com/d3v-friends/go-pure/fnReflect"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"testing"
	"time"
)
K
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

}