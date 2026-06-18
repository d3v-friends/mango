package tester

import (
	"time"

	"github.com/d3v-friends/go-tools/fnPointer"
	"github.com/d3v-friends/mango/v2/mgop"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ObjectIdArgs struct {
	Gte *bson.ObjectID `bson:"gte"`
	Eq  *bson.ObjectID `bson:"eq"`
}

func (x *ObjectIdArgs) Filter(filter bson.M, key string) bson.M {
	var compare = bson.M{}

	if !fnPointer.IsNil(x.Gte) {
		compare["$gte"] = *x.Gte
	}

	if !fnPointer.IsNil(x.Eq) {
		compare["$eq"] = *x.Eq
	}

	if len(compare) != 0 {
		filter[key] = compare
	}

	return filter
}

type Sorter string

const (
	SorterASC  Sorter = "ASC"
	SorterDESC Sorter = "DESC"
)

func (x Sorter) GetDirection() int32 {
	if x == SorterASC {
		return 1
	}
	return -1
}

type Pager struct {
	Page int64
	Size int64
}

func (x Pager) GetPage() int64 {
	return x.Page
}

func (x Pager) GetSize() int64 {
	return x.Size
}

func (x Pager) GetSkip() int64 {
	return x.Page * x.Size
}

func (x Pager) GetLimit() int64 {
	return x.Size
}

type StringArgs struct {
	Exact *string  `json:"exact,omitempty"`
	Like  *string  `json:"like,omitempty"`
	In    []string `json:"in,omitempty"`
}

func (x StringArgs) Filter(filter bson.M, key string) bson.M {
	if x.Exact != nil {
		filter[key] = *x.Exact
		return filter
	}

	if x.Like != nil {
		filter[key] = bson.M{
			mgop.Regex: *x.Like,
		}
		return filter
	}

	if 0 < len(x.In) {
		filter[key] = bson.M{
			mgop.In: x.In,
		}
		return filter
	}

	return filter
}

type Int64Args struct {
	Gt       *int64 `json:"gt,omitempty"`
	Gte      *int64 `json:"gte,omitempty"`
	Lt       *int64 `json:"lt,omitempty"`
	Lte      *int64 `json:"lte,omitempty"`
	Equal    *int64 `json:"equal,omitempty"`
	NotEqual *int64 `json:"notEqual,omitempty"`
}

func (x Int64Args) Filter(filter bson.M, key string) bson.M {
	var compare = bson.M{}

	if x.Gt != nil {
		compare[mgop.Gt] = *x.Gt
	}

	if x.Gte != nil {
		compare[mgop.Gte] = *x.Gte
	}

	if x.Lt != nil {
		compare[mgop.Lt] = *x.Lt
	}

	if x.Lte != nil {
		compare[mgop.Lte] = *x.Lte
	}

	if x.Equal != nil {
		compare[mgop.Eq] = *x.Equal
	}

	if x.NotEqual != nil {
		compare[mgop.Ne] = *x.NotEqual
	}

	if len(compare) == 0 {
		return filter
	}

	filter[key] = compare

	return filter
}

type TimeArgs struct {
	Gt       *time.Time `json:"gt,omitempty"`
	Gte      *time.Time `json:"gte,omitempty"`
	Lt       *time.Time `json:"lt,omitempty"`
	Lte      *time.Time `json:"lte,omitempty"`
	Equal    *time.Time `json:"equal,omitempty"`
	NotEqual *time.Time `json:"notEqual,omitempty"`
}

func (x TimeArgs) Filter(filter bson.M, key string) bson.M {
	var compare = bson.M{}

	if x.Gt != nil {
		compare[mgop.Gt] = *x.Gt
	}

	if x.Gte != nil {
		compare[mgop.Gte] = *x.Gte
	}

	if x.Lt != nil {
		compare[mgop.Lt] = *x.Lt
	}

	if x.Lte != nil {
		compare[mgop.Lte] = *x.Lte
	}

	if x.Equal != nil {
		compare[mgop.Eq] = *x.Equal
	}

	if x.NotEqual != nil {
		compare[mgop.Ne] = *x.NotEqual
	}

	if len(compare) == 0 {
		return filter
	}

	filter[key] = compare

	return filter
}
