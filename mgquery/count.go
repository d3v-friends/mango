package mgquery

import (
	"context"

	"github.com/d3v-friends/mango/v2"
	"github.com/d3v-friends/mango/v2/mgbuilder"
	"github.com/d3v-friends/mango/v2/mgctx"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Count[T mango.Model](
	ctx context.Context,
	filter any,
	opts ...*options.CountOptionsBuilder,
) (_ int64, err error) {
	var db *mongo.Database
	if db, err = mgctx.GetReaderDB(ctx); err != nil {
		return
	}

	var model = *new(T)
	var col = db.Collection(model.GetColNm())

	var opt = options.Count()
	if len(opts) == 1 {
		opt = opts[0]
	}

	var f = mgbuilder.Filter(filter)

	return col.CountDocuments(ctx, f, opt)
}
