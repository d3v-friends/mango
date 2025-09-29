package mgQuery

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/d3v-friends/go-tools/fnCase"
	"github.com/d3v-friends/go-tools/fnError"
	"github.com/d3v-friends/go-tools/fnPointer"
	"github.com/d3v-friends/mango/mgOp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
)

// 전제조건
// 1) mongodb field 는 반드시 lowerCamelCase 방식으로 작성 해야 한다.
// 2) 연속되는 대문자 불가능 -> 자동생성 방식으로 특정 불가능하기 때문
// 3) gql 모델과 mongodb 모델의 필드명이 완벽히 일치 해야함.

// 필드 추론 방식
// 1) 각 필드의 이름을 lowerCamelCase 로 변경
// 2) tag 가 존재한다면 그것을 우선한다.

// 필드는 다음 인터페이스 중 한가지를 반드시 구현한다
// StringArgs
// CompareArgs
// ArrayArgs

type Raw interface {
	Raw(registries ...*bsoncodec.Registry) (bson.Raw, error)
}

func ParseFilter(v any, registries ...*bsoncodec.Registry) (any, error) {
	switch t := v.(type) {
	case bson.M, bson.A, bson.D, bson.E, bson.Raw:
		return t, nil
	case Raw:
		return t.Raw(registries...)
	case nil:
		return bson.M{}, nil
	default:
		return parseFilterBsonM(bson.M{}, "", v)
	}
}

func ParseFilterBsonM(v any) (bson.M, error) {
	return parseFilterBsonM(bson.M{}, "", v)
}

func parseFilterBsonM(filter bson.M, parent string, v any) (_ bson.M, err error) {
	// gqlgen ID 로 바뀌는 것 수정
	parent = strings.ReplaceAll(parent, "ID", "Id")

	// mongodb 에서는 _id 강제로 사용
	if strings.ToLower(parent) == "id" {
		parent = "_id"
	}

	var vo = reflect.ValueOf(v)
	if vo.Kind() == reflect.Pointer && vo.IsNil() {
		return filter, nil
	}

	var f, isOk = v.(AppendFilterArgs)
	if isOk {
		return f.AppendFilter(filter, parent), nil
	}

	switch vo.Kind() {
	case reflect.Pointer:
		return parseFilterBsonM(filter, parent, vo.Elem().Interface())
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

			if filter, err = parseFilterBsonM(filter, key, field.Interface()); err != nil {
				return nil, err
			}
		}

		return filter, nil
	default:
		return filter, nil
	}
}

/* ------------------------------------------------------------------------------------------------------------ */

func ParseSorter(v any, registries ...*bsoncodec.Registry) (any, error) {
	if fnPointer.IsNil(v) {
		return bson.D{}, nil
	}

	switch t := v.(type) {
	case bson.M, bson.A, bson.D, bson.E:
		return t, nil
	case Raw:
		return t.Raw(registries...)
	case nil:
		return bson.D{}, nil
	default:
		return ParseSorterBsonD(v)
	}
}

func ParseSorterBsonD(v any) (sorter bson.D, err error) {
	sorter = bson.D{}
	var vo = reflect.ValueOf(v)
	switch vo.Kind() {
	case reflect.Slice:
		for i := 0; i < vo.Len(); i++ {
			var field = vo.Index(i)
			if !field.CanInterface() {
				continue
			}

			var elem *bson.E
			if elem, err = parserSorter("", field.Interface()); err != nil {
				return
			}
			sorter = append(sorter, *elem)
		}
	default:
		var elem *bson.E
		if elem, err = parserSorter("", v); err != nil {
			return
		}
		sorter = append(sorter, *elem)
	}
	return
}

const (
	ErrSorterIsNil   = "sorter_is_nil"
	ErrInvalidSorter = "invalid_sorter"
)

func parserSorter(parent string, v any) (res *bson.E, err error) {
	// gqlgen ID 로 바뀌는 것 수정
	parent = strings.ReplaceAll(parent, "ID", "Id")
	if strings.ToLower(parent) == "id" {
		parent = "_id"
	}

	var vo = reflect.ValueOf(v)
	if vo.Kind() == reflect.Pointer && vo.IsNil() {
		err = fnError.NewF(ErrSorterIsNil)
		return
	}

	var f, isOk = v.(SortArgs)
	if isOk {
		return &bson.E{
			Key:   parent,
			Value: f.GetDirection(),
		}, nil
	}

	switch vo.Kind() {
	case reflect.Pointer:
		return parserSorter(parent, vo.Elem().Interface())
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

			if res, err = parserSorter(key, field.Interface()); err != nil {
				err = nil
				continue
			}
			return
		}
		fallthrough
	default:
		err = fnError.NewF(ErrInvalidSorter)
		return
	}
}

/*------------------------------------------------------------------------------------------------*/

type AppendFilterArgs interface {
	AppendFilter(filter bson.M, key string) bson.M
}

// CompareArgs 존재하는 모든 조건이 적용된다. (and 연산)
type CompareArgs[T any] interface {
	GetGt() *T
	GetGte() *T
	GetLt() *T
	GetLte() *T
	GetEqual() *T
	GetNotEqual() *T
}

func AppendFilterCompareArgs[T any](
	filter bson.M,
	key string,
	args CompareArgs[T],
) bson.M {
	if fnPointer.IsNil(args) {
		return filter
	}

	var compare = bson.M{}
	if gt := args.GetGt(); !fnPointer.IsNil(gt) {
		compare[mgOp.Gt] = *gt
	}

	if gte := args.GetGte(); !fnPointer.IsNil(gte) {
		compare[mgOp.Gte] = *gte
	}

	if lt := args.GetLt(); !fnPointer.IsNil(lt) {
		compare[mgOp.Lt] = *lt
	}

	if lte := args.GetLte(); !fnPointer.IsNil(lte) {
		compare[mgOp.Lte] = *lte
	}

	if equal := args.GetEqual(); !fnPointer.IsNil(equal) {
		compare[mgOp.Eq] = *equal
	}

	if notEqual := args.GetNotEqual(); !fnPointer.IsNil(notEqual) {
		compare[mgOp.Ne] = *notEqual
	}

	if len(compare) == 0 {
		return filter
	}

	filter[key] = compare
	return filter
}

/* ------------------------------------------------------------------------------------------------------------ */

type PagerArgs interface {
	GetPage() int64
	GetSize() int64
	GetSkip() *int64
	GetLimit() *int64
}

type SortArgs interface {
	GetDirection() int32
}

func AppendSorter(
	sorter bson.D,
	key string,
	args SortArgs,
) bson.D {
	if fnPointer.IsNil(args) {
		return sorter
	}
	return append(sorter, bson.E{Key: key, Value: args.GetDirection()})
}

/* ------------------------------------------------------------------------------------------------------------ */

type ValueArgs[T any] interface {
	GetEqual() *T
	GetIn() []T
	GetHasAll() []T
	GetNin() []T
}

func AppendValueArgs[T any](filter bson.M, key string, args ValueArgs[T]) bson.M {
	if fnPointer.IsNil(args) {
		return filter
	}

	if equal := args.GetEqual(); !fnPointer.IsNil(equal) {
		filter[key] = *equal
		return filter
	}

	if in := args.GetIn(); 0 < len(in) {
		filter[key] = bson.M{
			mgOp.In: in,
		}
		return filter
	}

	if hasAll := args.GetHasAll(); 0 < len(hasAll) {
		filter[key] = bson.M{
			mgOp.All: hasAll,
		}
		return filter
	}

	if nin := args.GetNin(); 0 < len(nin) {
		filter[key] = bson.M{
			mgOp.Nin: nin,
		}
		return filter
	}

	return filter
}
