package mtx

import (
	"context"
	"github.com/d3v-friends/go-pure/fnEnv"
	"github.com/d3v-friends/go-pure/fnFile"
	"github.com/d3v-friends/mango/mctx"
	"github.com/d3v-friends/mango/mtype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type testModel struct {
	Id        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	CreatedAt time.Time          `bson:"createdAt"`
}

func (x *testModel) GetID() primitive.ObjectID {
	return x.Id
}

func (x *testModel) GetCollectionNm() string {
	return "testModel"
}

func (x *testModel) GetMigrateList() mtype.FnMigrateList {
	return mtype.FnMigrateList{}
}

type testTool struct {
	db *mongo.Database
}

func newTestTool(db *mongo.Database) (res *testTool) {
	res = &testTool{
		db: db,
	}
	return
}

func (x *testTool) Context() (ctx context.Context) {
	ctx = context.TODO()
	ctx = mctx.Set(ctx, x.db)
	return
}

func (x *testTool) ReadEnv(path fnFile.Path) (err error) {
	var p string
	if p, err = path.Path(); err != nil {
		return
	}

	if err = fnEnv.ReadFromFile(p); err != nil {
		return
	}

	return
}
