package mgBson

import (
	"github.com/d3v-friends/go-tools/fnPointer"
	"go.mongodb.org/mongo-driver/bson"
)

type SortArgs interface {
	GetDirection() int64
}

func AppendSorter(
	d bson.D,
	key string,
	value SortArgs,
) bson.D {
	if fnPointer.IsNil(value) {
		return d
	}
	return append(d, bson.E{
		Key:   key,
		Value: value.GetDirection(),
	})
}

func AppendD(
	d bson.D,
	key string,
	value any,
) bson.D {
	if fnPointer.IsNil(value) {
		return d
	}
	return append(d, bson.E{
		Key:   key,
		Value: value,
	})
}
