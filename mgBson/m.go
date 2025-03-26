package mgBson

import (
	"github.com/d3v-friends/go-tools/fnPointer"
	"go.mongodb.org/mongo-driver/bson"
)

func AppendM(
	m bson.M,
	key string,
	value any,
) bson.M {
	if fnPointer.IsNil(value) {
		return m
	}

	m[key] = value
	return m
}
