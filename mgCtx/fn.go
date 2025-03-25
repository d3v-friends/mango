package mgCtx

import (
	"context"
	"fmt"
	"github.com/d3v-friends/go-tools/fnCtx"
	"github.com/d3v-friends/go-tools/fnError"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

const ctxKeyMongoDB fnCtx.Key[*mongo.Database] = "CTX_MONGO_DATABASE"

const (
	ErrInvalidNameType = "invalid_name_type"
)

type Model interface {
	GetColNm() string
}

func SetDB(ctx context.Context, db *mongo.Database) context.Context {
	return fnCtx.Set(ctx, ctxKeyMongoDB, db)
}

func GetDB(ctx context.Context) (*mongo.Database, error) {
	return fnCtx.Get(ctx, ctxKeyMongoDB)
}

func GetDBP(ctx context.Context) *mongo.Database {
	return fnCtx.GetP(ctx, ctxKeyMongoDB)
}

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
	case Model:
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

func GetColByModel[T Model](ctx context.Context) (*mongo.Collection, error) {
	return GetCol(ctx, (*new(T)).GetColNm())
}

func GetColByModelP[T Model](ctx context.Context) *mongo.Collection {
	return GetColP(ctx, (*new(T)).GetColNm())
}
