package aggregate

import (
	"context"
	"fmt"
	"github.com/d3v-friends/mango/models"
	"github.com/d3v-friends/mango/mvars"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

type LookUpTag struct {
	From         string
	As           string
	LocalField   string
	ForeignField string
}

// todo 여기서부터 다시하기
func (x LookUpTag) Filter() bson.D {
	return bson.D{
		{
			Key: mvars.OLookUp,
			Value: bson.M{
				"from":         x.From,
				"as":           x.As,
				"localField":   x.LocalField,
				"foreignField": x.ForeignField,
			},
		},
		{
			Key: mvars.OMatch,
			Value: bson.M{
				"": bson.M{
					mvars.OExist: true,
					mvars.ONot: bson.M{
						mvars.OType: mvars.VArrayType,
					},
					mvars.OType: mvars.VObjectType,
				},
			},
		},
		{
			Key: mvars.OReplaceRoot,
			Value: bson.M{
				x.As: fmt.Sprintf("$%s", x.As),
			},
		},
	}
}

func LockUp[M any](
	ctx context.Context,
	db *mongo.Database,
	filter models.IfFilter,
) (err error) {
	var model = new(M)

	db.
		Collection(filter.GetCollectionNm()).
		Aggregate(ctx, bson.D{
			{
				Key:   mvars.OMatch,
				Value: filter.GetFilter(),
			},
			{
				Key:   mvars.OLookUp,
				Value: nil,
			},
		})

	panic("not impl")
}

func isStruct[M any]() (res bool) {
	var v = new(M)

	var to = reflect.TypeOf(v)
	if to.Kind() != reflect.Pointer {
		return
	}

	to = reflect.TypeOf(*v)
	if to.Kind() != reflect.Struct {
		return
	}

	res = true
	return
}
