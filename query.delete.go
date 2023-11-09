package mango

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/d3v-friends/go-pure/fnPanic"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DeleteOne(ctx context.Context, i IfFilter, opts ...*options.DeleteOptions) (err error) {
	var filter bson.M
	if filter, err = i.Filter(); err != nil {
		return
	}

	var col = GetMangoP(ctx).DB.Collection(i.ColNm())
	var count int64
	if count, err = col.CountDocuments(ctx, filter); err != nil {
		return
	}

	if count == 0 || 2 < count {
		err = fmt.Errorf("delete one must has 1 document: doc_count=%d, filter=%s", count, fnPanic.Get(json.Marshal(filter)))
		return
	}

	if _, err = col.DeleteOne(ctx, i, opts...); err != nil {
		return
	}

	return
}

func DeleteMany(ctx context.Context, i IfFilter, opts ...*options.DeleteOptions) (err error) {
	var filter bson.M
	if filter, err = i.Filter(); err != nil {
		return
	}

	var col = GetMangoP(ctx).DB.Collection(i.ColNm())
	if _, err = col.DeleteMany(ctx, filter, opts...); err != nil {
		return
	}

	return
}
