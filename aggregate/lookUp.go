package aggregate

import (
	"fmt"
	"github.com/d3v-friends/mango/mtype"
	"github.com/d3v-friends/mango/mvars"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"strings"
)

const (
	goTagLookUp = "lookUp"
	goTagBson   = "bson"
)

type tagLookUp struct {
	From         string
	As           string
	LocalField   string
	ForeignField string
	Type         mtype.ResType
}

func LookUp(v any) (res bson.A, err error) {
	res = make(bson.A, 0)

	var to = reflect.TypeOf(v)
	for i := 0; i < to.NumField(); i++ {
		var field = to.Field(i).Tag.Get(goTagLookUp)

		if field == "" {
			continue
		}

		var bsonAs = to.Field(i).Tag.Get(goTagBson)
		if bsonAs == "" {
			err = fmt.Errorf("not found bson tag: model=%s", to.Name())
			return
		}

		var tag *tagLookUp
		if tag, err = newTagLookUp(bsonAs, field); err != nil {
			return
		}

		var kind = to.Field(i).Type.Kind()
		switch kind {
		case reflect.Pointer:
			tag.Type = mtype.ResTypeObject
		case reflect.Array:
			tag.Type = mtype.ResTypeObject
		default:
			err = fmt.Errorf("invalid field type: kind=%s", kind.String())
			return
		}

		if err = tag.valid(); err != nil {
			return
		}

		res = append(res, bson.M{
			mvars.OLookUp: bson.M{
				"from":         tag.From,
				"as":           tag.As,
				"localField":   tag.LocalField,
				"foreignField": tag.ForeignField,
			},
		})

		if tag.Type == mtype.ResTypeObject {
			res = append(res, bson.M{
				mvars.OAddField: bson.M{
					tag.As: bson.M{
						mvars.OArrayElemAt: bson.A{
							fmt.Sprintf("$%s", tag.As), 0,
						},
					},
				},
			})
		}
	}

	return
}

func newTagLookUp(as, v string) (res *tagLookUp, err error) {
	res = &tagLookUp{
		As: as,
	}

	defer func() {
		if err != nil {
			res = nil
		}
	}()

	valueList := strings.Split(v, ";")
	for _, value := range valueList {
		itemList := strings.Split(value, ":")
		if len(itemList) != 2 {
			err = fmt.Errorf("invalid tag value: tag=%s", value)
			return
		}

		var tagKey, tagValue = itemList[0], itemList[1]
		switch tagKey {
		case "from":
			res.From = tagValue
		case "localField":
			res.LocalField = tagValue
		case "foreignField":
			res.ForeignField = tagValue
		default:
			err = fmt.Errorf("invalid tag key: tag=%s", value)
			return
		}
	}

	if err = res.valid(); err != nil {
		return
	}

	return
}

func (x tagLookUp) valid() (err error) {
	if x.As == "" {
		err = fmt.Errorf("not found lookup 'as' tag value")
		return
	}

	if x.From == "" {
		err = fmt.Errorf("not found lookup 'from' tag value")
		return
	}

	if x.ForeignField == "" {
		err = fmt.Errorf("not found lookup 'foreignField' tag value")
		return
	}

	if x.LocalField == "" {
		err = fmt.Errorf("not found lookup 'localField' tag value")
		return
	}

	return
}
