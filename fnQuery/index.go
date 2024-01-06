package fnQuery

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type Time struct {
	Equal    *time.Time
	NotEqual *time.Time
	GT       *time.Time
	GTE      *time.Time
	LT       *time.Time
	LTE      *time.Time
}

func (x *Time) Query(nm string) (res bson.M, err error) {
	var query = make(bson.M)
	if x.Equal != nil {
		query["$equal"] = *x.Equal
	}

	if x.NotEqual != nil {
		query["$ne"] = *x.NotEqual
	}

	if x.GT != nil {
		query["$gt"] = *x.GT
	}

	if x.GTE != nil {
		query["$gte"] = *x.GTE
	}

	if x.LT != nil {
		query["$lt"] = *x.LT
	}

	if x.LTE != nil {
		query["$lte"] = *x.LTE
	}

	if len(query) == 0 {
		err = fmt.Errorf("not found query value")
		return
	}

	res = bson.M{
		nm: query,
	}

	return
}

/*------------------------------------------------------------------------------------------------*/

type NumberOperator string

const (
	NumberOperatorEqual    NumberOperator = "$equal"
	NumberOperatorNotEqual NumberOperator = "$ne"
	NumberOperatorGT       NumberOperator = "$gt"
	NumberOperatorGTE      NumberOperator = "$gte"
	NumberOperatorLT       NumberOperator = "$lt"
	NumberOperatorLTE      NumberOperator = "$lte"
)

func (x *NumberOperator) String() string {
	return string(*x)
}

/*------------------------------------------------------------------------------------------------*/

type String struct {
	Value    string
	Operator StringOperator
}

type StringOperator string

const (
	OperatorStringEqual StringOperator = "EQUAL"
	OperatorStringLike  StringOperator = "LIKE"
)

type IdxDir int64

const (
	IdxDirASC  IdxDir = 1
	IdxDirDESC IdxDir = -1
)

func (x *IdxDir) ToInt() int {
	return int(*x)
}
