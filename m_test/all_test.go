package m_test

import (
	"context"
	"github.com/d3v-friends/go-pure/fnEnv"
	"github.com/d3v-friends/go-pure/fnPanic"
	"github.com/d3v-friends/go-pure/fnParams"
	"github.com/d3v-friends/mango"
	"github.com/d3v-friends/mango/m_codec"
	"github.com/d3v-friends/mango/m_ctx"
	"github.com/d3v-friends/mango/m_migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type DocTest struct {
	Id        primitive.ObjectID `bson:"_id"`
	IsTx      bool               `bson:"isTx"`
	Content   string             `bson:"content"`
	CreatedAt time.Time          `bson:"createdAt"`
}

func (x *DocTest) GetID() primitive.ObjectID {
	return x.Id
}

func (x *DocTest) GetColNm() string {
	return docTestNm
}

func (x *DocTest) GetTxNm() string {
	return "isTx"
}

func (x *DocTest) GetMigrateList() m_migrate.FnMigrateList {
	return mgDocTest
}

const docTestNm = "docTests"

var mgDocTest = m_migrate.FnMigrateList{
	func(ctx context.Context, col *mongo.Collection) (memo string, err error) {
		memo = "indexing content"
		_, err = col.InsertOne(ctx, mongo.IndexModel{
			Keys: bson.D{
				{
					Key:   "content",
					Value: 1,
				},
			},
		})
		return
	},
}

/* ------------------------------------------------------------------------------------------------------------ */
type TestTool struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func (x *TestTool) Context() (ctx context.Context) {
	ctx = context.TODO()
	ctx = m_ctx.SetDB(ctx, x.DB)
	return
}

func NewTestTool(truncate ...bool) (res *TestTool) {
	fnPanic.On(fnEnv.ReadFromFile("../env/.env"))

	res = &TestTool{}
	res.Client = fnPanic.Get(mango.NewClient(&mango.IConn{
		Host:        fnEnv.Read("MG_HOST"),
		Username:    fnEnv.Read("MG_USERNAME"),
		Password:    fnEnv.Read("MG_PASSWORD"),
		SetRegistry: m_codec.RegisterDecimal,
	}))

	res.DB = res.Client.Database(fnEnv.Read("MG_DATABASE"))

	if fnParams.Get(truncate) {
		fnPanic.On(res.DB.Drop(context.TODO()))
	}

	return
}
