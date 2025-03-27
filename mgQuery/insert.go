package mgQuery

import (
	"context"
	"github.com/d3v-friends/go-tools/fnError"
	"github.com/d3v-friends/mango"
	"github.com/d3v-friends/mango/mgCtx"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ErrEmptyModels = "empty_models"
)

// InsertOne
// T 제네릭을 명시하지 않으면 new(T) 의 T 를 정확하게 인식 못한다. (golang 의 문제, 또는 golang 의 제네릭 이해의 문제)
// 다른 query helper 들은 명시적으로 입력해야 하는 만큼 제네릭타입의 new 가 가능함
func InsertOne[T mango.Model](
	ctx context.Context,
	model T,
) (err error) {
	var col *mongo.Collection

	if col, err = mgCtx.GetCol(ctx, model); err != nil {
		return
	}

	if _, err = col.InsertOne(ctx, model); err != nil {
		return
	}

	return
}

func InsertMany[T mango.Model](
	ctx context.Context,
	models []T,
) (err error) {
	if len(models) == 0 {
		err = fnError.New(ErrEmptyModels)
		return
	}

	var ls = make([]interface{}, len(models))

	for i, model := range models {
		ls[i] = model
	}

	var col *mongo.Collection
	if col, err = mgCtx.GetCol(ctx, models[0]); err != nil {
		return
	}

	if _, err = col.InsertMany(ctx, ls); err != nil {
		return
	}

	return
}
