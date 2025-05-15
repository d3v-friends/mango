package mgQuery

import (
	"fmt"
	"github.com/d3v-friends/go-tools/fnCase"
	"github.com/d3v-friends/go-tools/fnPointer"
	"github.com/d3v-friends/mango/mgOp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"reflect"
	"strings"
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
		return parseSorterBsonD(bson.D{}, "", v)
	}
}

func ParseSorterBsonD(v any) (bson.D, error) {
	return parseSorterBsonD(bson.D{}, "", v)
}

func parseSorterBsonD(sorter bson.D, parent string, v any) (_ bson.D, err error) {
	// gqlgen ID 로 바뀌는 것 수정
	parent = strings.ReplaceAll(parent, "ID", "Id")

	if strings.ToLower(parent) == "id" {
		parent = "_id"
	}

	var vo = reflect.ValueOf(v)
	if vo.Kind() == reflect.Pointer && vo.IsNil() {
		return sorter, nil
	}

	var f, isOk = v.(SortArgs)
	if isOk {
		return AppendSorter(sorter, parent, f), nil
	}

	switch vo.Kind() {
	case reflect.Pointer:
		return parseSorterBsonD(sorter, parent, vo.Elem().Interface())
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

			if sorter, err = parseSorterBsonD(sorter, key, field.Interface()); err != nil {
				return
			}
		}
		return sorter, nil
	default:
		return sorter, nil
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

// StringArgs 3개의 조건중 1개만 적용된다.
// 순서 중요. 우선 적용은
// 1) exact
// 2) like
// 3) in
// 순서대로 제일 먼저 있는 값을 적용한다.
type StringArgs interface {
	GetExact() *string
	GetLike() *string
	GetIn() []string
}

func AppendFilterStringArgs(
	filter bson.M,
	key string,
	args StringArgs,
) bson.M {
	if fnPointer.IsNil(args) {
		return filter
	}

	if exact := args.GetExact(); !fnPointer.IsNil(exact) {
		filter[key] = *exact
		return filter
	}

	if like := args.GetLike(); !fnPointer.IsNil(like) {
		filter[key] = bson.M{
			mgOp.Regex: *like,
		}
		return filter
	}

	if in := args.GetIn(); len(in) != 0 {
		filter[key] = bson.M{
			mgOp.In: in,
		}
		return filter
	}

	return filter
}

type PagerArgs interface {
	GetPage() int64
	GetSize() int64
	GetSkip() *int64
	GetLimit() *int64
}

type SortArgs interface {
	GetDirection() int64
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

// ArrayArgs 3개의 조건중 1개만 적용된다.
// 순서 중요. 우선 적용은
// 1) equal
// 2) hasAll
// 3) in
// 순서대로 제일 먼저 있는 값을 적용한다.

type ArrayArgs[T any] interface {
	GetEqual() *T
	GetIn() []T
	GetHasAll() []T
}

func AppendFilterArrayArgs[T any](
	filter bson.M,
	key string,
	args ArrayArgs[T],
) bson.M {
	if fnPointer.IsNil(args) {
		return filter
	}

	if equal := args.GetEqual(); !fnPointer.IsNil(equal) {
		filter[key] = equal
		return filter
	}

	var hasAll = args.GetHasAll()
	if len(hasAll) != 0 {
		filter[key] = hasAll
		return filter
	}

	var in = args.GetIn()
	if len(in) != 0 {
		filter[key] = bson.M{
			mgOp.In: in,
		}
		return filter
	}

	return filter
}
