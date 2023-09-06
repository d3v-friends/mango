package mctx

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestCtx(testing *testing.T) {
	ctx := context.TODO()
	var db = GetP[*mongo.Database](ctx)
	db.Client()
}
