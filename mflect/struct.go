package mflect

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
)

type RequiredField struct {
	BsonNm string
	Type   reflect.Type
}

var requiredFields = []*RequiredField{
	{
		BsonNm: "_id",
		Type:   reflect.TypeOf(primitive.ObjectID{}),
	},
	{
		BsonNm: "isLock",
		Type:   reflect.TypeOf(primitive.ObjectID{}),
	},
}

func IsModel[T any](v T) (err error) {
	panic("not impl")
}
