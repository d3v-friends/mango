package mFilter

import (
	"github.com/d3v-friends/mango"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

func NewId(
	colNm string,
	id primitive.ObjectID,
) mango.IfFilter {
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

func NewIds(
	colNm string,
	operator Operator,
	ids ...primitive.ObjectID,
) mango.IfFilter {
	return &Ids{
		ids:      ids,
		operator: operator,
		colNm:    colNm,
	}
}
