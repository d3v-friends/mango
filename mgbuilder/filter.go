package mgbuilder

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/d3v-friends/go-tools/fnCase"
	"github.com/d3v-friends/go-tools/fnPointer"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func Filter(ctx context.Context, v any) bson.M {
	switch t := v.(type) {
	case bson.M:
		return t
	case nil:
		return bson.M{}
	default:
		return convertToBsonM(
			ctx,
			bson.M{},
			"",
			v,
		)
	}
}

type FilterArg interface {
	Filter(ctx context.Context) any
}

func convertToBsonM(
	ctx context.Context,
	filter bson.M,
	parent string,
	v any,
) (_ bson.M) {
	if fnPointer.IsNil(v) {
		return filter
	}

	// gqlgen ID 로 바뀌는 것 수정
	parent = strings.ReplaceAll(parent, "ID", "Id")

	// mongodb 에서는 _id 강제로 사용
	if strings.ToLower(parent) == "id" {
		parent = "_id"
	}

	var vo = reflect.ValueOf(v)

	var f, isOk = v.(FilterArg)
	if isOk {
		var value = f.Filter(ctx)
		if fnPointer.IsNotNil(value) {
			filter[parent] = value
		}
		return filter
	}

	switch vo.Kind() {
	case reflect.Pointer:
		return convertToBsonM(ctx, filter, parent, vo.Elem().Interface())
	case reflect.Struct:
		for i := 0; i < vo.NumField(); i++ {
			var field = vo.Field(i)
			if !field.CanInterface() {
				continue
			}

			var key = fnCase.CamelCase(reflect.TypeOf(v).Field(i).Name)
			if parent != "" {
				key = fmt.Sprintf("%s.%s", parent, key)
			}

			filter = convertToBsonM(ctx, filter, key, field.Interface())
		}

		return filter
	default:
		return filter
	}
}
