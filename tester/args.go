package tester

import (
	"context"
	"time"

	"github.com/d3v-friends/go-tools/fnPointer"
	"github.com/d3v-friends/mango/v2/mgop"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ObjectIdArgs struct {
	Gte *bson.ObjectID
	Eq  *bson.ObjectID
}

func (x *ObjectIdArgs) Filter(_ context.Context) any {
	var compare = bson.M{}

	if fnPointer.IsNotNil(x.Gte) {
		compare[mgop.Gte] = *x.Gte
	}

	if fnPointer.IsNotNil(x.Eq) {
		compare[mgop.Eq] = *x.Eq
	}

	if len(compare) == 0 {
		return nil
	}

	return compare
}

type Sorter string

const (
	SorterASC  Sorter = "ASC"
	SorterDESC Sorter = "DESC"
)

func (x Sorter) GetDirection(_ context.Context) int32 {
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

func (x StringArgs) Filter(_ context.Context) any {
	if x.Exact != nil {
		return *x.Exact
	}

	if x.Like != nil {
		return bson.M{
			mgop.Regex: *x.Like,
		}
	}

	if 0 < len(x.In) {
		return bson.M{
			mgop.In: x.In,
		}

	}

	return nil
}

type Int64Args struct {
	Gt       *int64 `json:"gt,omitempty"`
	Gte      *int64 `json:"gte,omitempty"`
	Lt       *int64 `json:"lt,omitempty"`
	Lte      *int64 `json:"lte,omitempty"`
	Equal    *int64 `json:"equal,omitempty"`
	NotEqual *int64 `json:"notEqual,omitempty"`
}

func (x Int64Args) Filter(_ context.Context) any {
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
		return nil
	}

	return compare
}

type TimeArgs struct {
	Gt       *time.Time `json:"gt,omitempty"`
	Gte      *time.Time `json:"gte,omitempty"`
	Lt       *time.Time `json:"lt,omitempty"`
	Lte      *time.Time `json:"lte,omitempty"`
	Equal    *time.Time `json:"equal,omitempty"`
	NotEqual *time.Time `json:"notEqual,omitempty"`
}

func (x TimeArgs) Filter(_ context.Context) any {
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
		return nil
	}

	return compare
}

type DecimalArgs struct {
	Gt       *decimal.Decimal `json:"gt,omitempty"`
	Gte      *decimal.Decimal `json:"gte,omitempty"`
	Lt       *decimal.Decimal `json:"lt,omitempty"`
	Lte      *decimal.Decimal `json:"lte,omitempty"`
	Equal    *decimal.Decimal `json:"equal,omitempty"`
	NotEqual *decimal.Decimal `json:"notEqual,omitempty"`
}

func (x DecimalArgs) Filter(_ context.Context) any {
	var compare = bson.M{}

	if fnPointer.IsNotNil(x.Gt) {
		compare[mgop.Gt] = *x.Gt
	}

	if fnPointer.IsNotNil(x.Gte) {
		compare[mgop.Gte] = *x.Gte
	}

	if fnPointer.IsNotNil(x.Lt) {
		compare[mgop.Lt] = *x.Lt
	}

	if fnPointer.IsNotNil(x.Lte) {
		compare[mgop.Lte] = *x.Lte
	}

	if fnPointer.IsNotNil(x.Equal) {
		compare[mgop.Eq] = *x.Equal
	}

	if fnPointer.IsNotNil(x.NotEqual) {
		compare[mgop.Ne] = *x.NotEqual
	}

	if len(compare) == 0 {
		return nil
	}

	return compare
}
