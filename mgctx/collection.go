package mgctx

import (
	"context"
	"fmt"
	"reflect"

	"github.com/d3v-friends/go-tools/fnCtx"
	"github.com/d3v-friends/go-tools/fnError"
	"github.com/d3v-friends/mango/v2"
	"github.com/d3v-friends/mango/v2/mgvalue"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	ctxKeyMongoDB       fnCtx.Key[*mongo.Database] = "CTX_MONGO_DATABASE"
	ctxKeyReaderMongoDB fnCtx.Key[*mongo.Database] = "CTX_READER_MONGO_DATABASE"
	ctxKeyWriterMongoDB fnCtx.Key[*mongo.Database] = "CTX_WRITER_MONGO_DATABASE"
)

func GetReaderCollection(ctx context.Context, name any) (col *mongo.Collection, err error) {
	var db *mongo.Database
	if db, err = GetReaderDB(ctx); err != nil {
		return
	}

	var colName string
	if colName, err = getCollectionName(name); err != nil {
		return
	}

	col = db.Collection(colName)
	return
}

func GetWriterCollection(ctx context.Context, name any) (col *mongo.Collection, err error) {
	var db *mongo.Database
	if db, err = GetWriterDB(ctx); err != nil {
		return
	}

	var colName string
	if colName, err = getCollectionName(name); err != nil {
		return
	}

	col = db.Collection(colName)
	return
}

func getCollectionName(name any) (res string, err error) {
	switch t := name.(type) {
	case string:
		res = t
		return
	case *string:
		if !reflect.ValueOf(t).CanInterface() {
			err = fnError.New(mgvalue.ErrInvalidNameType)
			return
		}
		res = *t
		return
	case fmt.Stringer:
		res = t.String()
		return
	case mango.Model:
		res = t.GetColNm()
		return
	default:
		err = fnError.New(mgvalue.ErrInvalidNameType)
		return
	}
}
