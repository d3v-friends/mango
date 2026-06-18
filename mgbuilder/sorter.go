package mgbuilder

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/d3v-friends/go-tools/fnCase"
	"github.com/d3v-friends/go-tools/fnPointer"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func Sorter(v any) bson.D {
	switch t := v.(type) {
	case bson.D:
		return t
	case nil:
		return bson.D{}
	default:
		return BsonD(v)
	}
}

func BsonD(v any) (sorter bson.D) {
	sorter = bson.D{}
	var vo = reflect.ValueOf(v)
	switch vo.Kind() {
	case reflect.Slice:
		for i := 0; i < vo.Len(); i++ {
			var field = vo.Index(i)
			if !field.CanInterface() {
				continue
			}

			var elem bson.E
			if elem = convertToBsonE("", field.Interface()); elem.Key == "" {
				continue
			}
			sorter = append(sorter, elem)
		}
	default:
		var elem bson.E
		if elem = convertToBsonE("", v); elem.Key == "" {
			return
		}
		sorter = append(sorter, elem)
	}
	return
}

type SortArgs interface {
	GetDirection() int32
}

func convertToBsonE(parent string, v any) (res bson.E) {
	res = bson.E{}

	if fnPointer.IsNil(v) {
		return
	}

	// gqlgen ID 로 바뀌는 것 수정
	parent = strings.ReplaceAll(parent, "ID", "Id")
	if strings.ToLower(parent) == "id" {
		parent = "_id"
	}

	var vo = reflect.ValueOf(v)

	var f, isOk = v.(SortArgs)
	if isOk {
		res = bson.E{
			Key:   parent,
			Value: f.GetDirection(),
		}
		return
	}

	switch vo.Kind() {
	case reflect.Pointer:
		return convertToBsonE(parent, vo.Elem().Interface())
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

			if res = convertToBsonE(key, field.Interface()); res.Key == "" {
				continue
			}
			return
		}
		fallthrough
	default:
		return
	}
}
