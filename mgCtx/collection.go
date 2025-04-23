package mgCtx

import (
	"context"
	"fmt"
	"github.com/d3v-friends/go-tools/fnCtx"
	"github.com/d3v-friends/go-tools/fnError"
	"github.com/d3v-friends/mango"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

const (
	ctxKeyMongoDB      fnCtx.Key[*mongo.Database] = "CTX_MONGO_DATABASE"
	ErrInvalidNameType                            = "invalid_name_type"
)

func GetCol(ctx context.Context, name any) (col *mongo.Collection, err error) {
	var db *mongo.Database
	if db, err = GetDB(ctx); err != nil {
		return
	}

	switch t := name.(type) {
	case string:
		col = db.Collection(t)
		return
	case *string:
		if reflect.ValueOf(t).CanInterface() {
			col = db.Collection(*t)
			return
		} else {
			err = fnError.NewFields(ErrInvalidNameType, map[string]any{
				"string": nil,
			})
			return
		}
	case fmt.Stringer:
		col = db.Collection(t.String())
		return
	case mango.Model:
		col = db.Collection(t.GetColNm())
		return
	default:
		err = fnError.New(ErrInvalidNameType)
		return
	}
}

func GetColP(ctx context.Context, name any) *mongo.Collection {
	var col, err = GetCol(ctx, name)
	if err != nil {
		panic(err)
	}
	return col
}
