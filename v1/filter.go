package v1

import (
	"encoding/json"
	"fmt"
	"github.com/d3v-friends/go-pure/fnPanic"
	"github.com/d3v-friends/go-pure/fnParams"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strings"
	"time"
)

const (
	tagMango     = "mango"
	unitOperator = "operator"
)

func pickTag(tag string, key string) (value string, err error) {
	tag = strings.ReplaceAll(tag, " ", "")
	var tagLs = strings.Split(tag, ";")

	for _, unit := range tagLs {
		if !strings.HasPrefix(unit, key) {
			continue
		}

		var unitLs = strings.Split(unit, ":")
		if len(unitLs) != 2 {
			err = fmt.Errorf("mango tag is not valid: tag=%s, unit=%s", tag, unit)
			return
		}

		value = unitLs[1]
		return
	}

	err = fmt.Errorf("not found mango tag: tag=%s, unit_key=%s", tag, key)
	return
}

/* ------------------------------------------------------------------------------------------------------------ */

type Operator string

func (x Operator) String() string {
	return string(x)
}

type Operators []Operator

func (x Operators) Has(operator any) (err error) {
	var value = Operator(fnParams.ToStringP(operator))
	for _, oper := range x {
		if oper == value {
			return
		}
	}

	err = fmt.Errorf("invalid operator: operator=%s, operator_set=%s",
		value.String(),
		fnPanic.Get(json.Marshal(x)),
	)
	return
}

const (
	OperatorIn    Operator = "$in"
	OperatorGt    Operator = "$gt"
	OperatorGte   Operator = "$gte"
	OperatorLt    Operator = "$lt"
	OperatorLte   Operator = "$lte"
	OperatorEqual Operator = "$equal"
	OperatorNe    Operator = "$ne"
)

var operatorTime = Operators{
	OperatorGt,
	OperatorGte,
	OperatorLt,
	OperatorLte,
	OperatorEqual,
	OperatorNe,
}

/* ------------------------------------------------------------------------------------------------------------ */

// filter

type Id struct {
	id    primitive.ObjectID
	colNm string
}

func (x Id) Filter() (res bson.M, _ error) {
	res = bson.M{
		"_id": x.id,
	}
	return
}

func (x Id) ColNm() string {
	return x.colNm
}

func NewIdFilter(
	colNm string,
	id primitive.ObjectID,
) IfFilter {
	return &Id{
		id:    id,
		colNm: colNm,
	}
}

/* ------------------------------------------------------------------------------------------------------------ */

// Ids In operation
type Ids struct {
	ids      []primitive.ObjectID
	operator Operator
	colNm    string
}

func (i Ids) Filter() (res bson.M, _ error) {
	res = bson.M{
		i.operator.String(): bson.M{
			"_id": i.ids,
		},
	}
	return
}

func (i Ids) ColNm() string {
	return i.colNm
}

func NewIdsFilter(
	colNm string,
	operator Operator,
	ids ...primitive.ObjectID,
) IfFilter {
	return &Ids{
		ids:      ids,
		operator: operator,
		colNm:    colNm,
	}
}

/* ------------------------------------------------------------------------------------------------------------ */

type Period struct {
	GT       *time.Time `mango:"operator:$gt"`
	GTE      *time.Time `mango:"operator:$gte"`
	LT       *time.Time `mango:"operator:$lt"`
	LTE      *time.Time `mango:"operator:$lte"`
	Equal    *time.Time `mango:"operator:$equal"`
	NotEqual *time.Time `mango:"operator:$ne"`
}

func (x *Period) Filter() (filter bson.M, err error) {

	filter = make(bson.M)
	var valOf = reflect.ValueOf(*x)
	var typOf = reflect.TypeOf(*x)

	var fieldMax = typOf.NumField()
	for i := 0; i < fieldMax; i++ {
		if typOf.Field(i).Type.String() != "*time.Time" {
			continue
		}

		if valOf.Field(i).IsNil() {
			continue
		}

		var tag = typOf.Field(i).Tag.Get(tagMango)
		if tag == "" {
			continue
		}

		var operator string
		if operator, err = pickTag(tag, unitOperator); err != nil {
			return
		}

		if err = operatorTime.Has(operator); err != nil {
			return
		}

		filter[operator] = valOf.Field(i).Elem().Interface()
	}

	if len(filter) == 0 {
		err = fmt.Errorf("period is empty value")
		return
	}

	return
}

type IfPeriod interface {
	GT() *time.Time
	GTE() *time.Time
	LT() *time.Time
	LTE() *time.Time
	Equal() *time.Time
	NotEqual() *time.Time
}

func TimeFilter(v IfPeriod) (bson.M, error) {
	var period = &Period{
		GT:       v.GT(),
		GTE:      v.GTE(),
		LT:       v.LT(),
		LTE:      v.LTE(),
		Equal:    v.Equal(),
		NotEqual: v.NotEqual(),
	}

	return period.Filter()
}
