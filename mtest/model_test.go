package mtest

import (
	"context"
	"github.com/d3v-friends/go-pure/fnReflect"
	"github.com/d3v-friends/mango/mtype"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var modelAll = []mtype.IfMigrateModel{
	&Account{},
	&AccountData{},
}

type UserType string

const (
	UserTypeGeneral UserType = "General"
	UserTypeAdmin   UserType = "Admin"
)

type Account struct {
	Id            primitive.ObjectID `bson:"_id"`
	AccountDataId primitive.ObjectID `bson:"accountDataId"`
	UserType      UserType           `bson:"userType"`
	UpdatedAt     time.Time          `bson:"createdAt"`
}

func (x Account) GetID() primitive.ObjectID {
	return x.Id
}

const colAccount = "accounts"

var mgAccount = mtype.FnMigrateList{
	func(ctx context.Context, collection *mongo.Collection) (migrationNm string, err error) {
		migrationNm = "init indexing"
		if _, err = collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
			{
				Keys: bson.D{
					{
						Key:   "accountDataId",
						Value: 1,
					},
				},
				Options: &options.IndexOptions{
					Unique: fnReflect.ToPointer(true),
				},
			},
		}); err != nil {
			return
		}
		return
	},
}

func (x Account) GetCollectionNm() string {
	return colAccount
}

func (x Account) GetMigrateList() mtype.FnMigrateList {
	return mgAccount
}

type AccountData struct {
	Id        primitive.ObjectID `bson:"_id"`
	AccountId primitive.ObjectID `bson:"accountId"`
	Name      string             `bson:"name"`
	CreatedAt time.Time          `bson:"createdAt"`
}

const colAccountData = "accountData"

var mgAccountData = mtype.FnMigrateList{
	func(ctx context.Context, collection *mongo.Collection) (migrationNm string, err error) {
		migrationNm = "init indexing"
		_, err = collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
			{
				Keys: bson.D{
					{
						Key:   "accountId",
						Value: 1,
					},
				},
			},
		})
		return
	},
}

func (x AccountData) GetCollectionNm() string {
	return colAccountData
}

func (x AccountData) GetMigrateList() mtype.FnMigrateList {
	return mgAccountData
}

func (x AccountData) GetID() primitive.ObjectID {
	return x.Id
}

type AccountWithData struct {
	Account *Account     `bson:"inline"`
	Data    *AccountData `bson:"data" lookUp:"from:accountData;localField:accountDataId;foreignField:_id"`
}

func (x AccountWithData) GetCollectionNm() string {
	return colAccount
}
