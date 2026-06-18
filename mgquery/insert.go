package mgquery

import (
	"context"

	"github.com/d3v-friends/go-tools/fnError"
	"github.com/d3v-friends/mango/v2"
	"github.com/d3v-friends/mango/v2/mgctx"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	ErrEmptyModels = "empty_models"
)

// InsertOne
// T 제네릭을 명시하지 않으면 new(T) 의 T 를 정확하게 인식 못한다. (golang 의 문제, 또는 golang 의 제네릭 이해의 문제)
// 다른 query helper 들은 명시적으로 입력해야 제네릭타입의 new 가 가능함
func InsertOne[T mango.Model](
	ctx context.Context,
	model T,
	opts ...*options.InsertOneOptionsBuilder,
) (_ *mongo.InsertOneResult, err error) {
	var col *mongo.Collection

	if col, err = mgctx.GetCol(ctx, model); err != nil {
		return
	}

	var opt = options.InsertOne()
	if len(opts) == 1 {
		opt = opts[0]
	}

	return col.InsertOne(ctx, model, opt)
}

func InsertMany[T mango.Model](
	ctx context.Context,
	models []T,
	opts ...*options.InsertManyOptionsBuilder,
) (_ *mongo.InsertManyResult, err error) {
	if len(models) == 0 {
		err = fnError.New(ErrEmptyModels)
		return
	}

	var ls = make([]interface{}, len(models))

	for i, model := range models {
		ls[i] = model
	}

	var col *mongo.Collection
	if col, err = mgctx.GetCol(ctx, models[0]); err != nil {
		return
	}

	var opt = options.InsertMany()
	if len(opts) == 1 {
		opt = opts[0]
	}

	return col.InsertMany(ctx, ls, opt)
}
