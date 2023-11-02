package mTest

import (
	"context"
	"github.com/d3v-friends/go-pure/fnEnv"
	"github.com/d3v-friends/go-pure/fnPanic"
	"github.com/d3v-friends/go-pure/fnParams"
	"github.com/d3v-friends/mango"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type DocTest struct {
	Id        primitive.ObjectID `bson:"_id"`
	GroupId   primitive.ObjectID `bson:"groupId"`
	Name      string             `bson:"name"`
	CreatedAt time.Time          `bson:"createdAt"`
}

func (x *DocTest) GetID() primitive.ObjectID {
	return x.Id
}

func (x *DocTest) GetColNm() string {
	return docTestNm
}

func (x *DocTest) GetMigrateList() mMigrate.FnMigrateList {
	return mgDocTest
}

const docTestNm = "docTests"

var mgDocTest = mMigrate.FnMigrateList{
	func(ctx context.Context, col *mongo.Collection) (memo string, err error) {
		memo = "indexing name"
		_, err = col.Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys: bson.D{
				{
					Key:   "name",
					Value: 1,
				},
			},
		})
		return
	},
}

/* ------------------------------------------------------------------------------------------------------------ */

func NewMango(truncate ...bool) (res *mango.Mango) {
	fnPanic.On(fnEnv.ReadFromFile("../env/.env"))
	res = fnPanic.Get(mango.NewMango(
		&mango.IConn{
			Host:        fnEnv.Read("MG_HOST"),
			Username:    fnEnv.Read("MG_USERNAME"),
			Password:    fnEnv.Read("MG_PASSWORD"),
			Database:    fnEnv.Read("MG_DATABASE"),
			SetRegistry: mcodec.RegisterDecimal,
		},
	))

	if fnParams.Get(truncate) {
		fnPanic.On(res.Truncate(context.TODO()))
	}

	return
}
