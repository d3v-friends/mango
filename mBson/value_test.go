package mBson_test

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestUpdateBson(test *testing.T) {

	test.Run("kind vs value", func(t *testing.T) {
		var v = &Test{
			Name: "hello",
		}

		var typeOf = reflect.TypeOf(v)

		assert.Equal(t, typeOf.Kind(), reflect.Pointer)
		assert.Equal(t, typeOf.Elem().Kind(), reflect.Struct)

	})
}

type Test struct {
	Name string `bson:"name"`
}
