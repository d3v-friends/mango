package mango

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type Mango struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMango(ctx context.Context, i *IConn) (res *Mango, err error) {
	if i == nil {
		err = fmt.Errorf("not found mango.IConn")
		return
	}

	res = &Mango{}
	if res.client, err = mongo.Connect(ctx, i.options()); err != nil {
		return
	}

	if err = res.client.Ping(ctx, readpref.Primary()); err != nil {
		return
	}

	res.db = res.client.Database(i.Database)
	return
}

func (x *Mango) Database() *mongo.Database {
	return x.db
}

func (x *Mango) Migrate(ctx context.Context) (err error) {
	panic("not impl")
}

/* ------------------------------------------------------------------------------------------------------------ */

type FnSetRegistry func(registry *bsoncodec.Registry) *bsoncodec.Registry

type IConn struct {
	Host        string
	Username    string
	Password    string
	Database    string
	SetRegistry FnSetRegistry
}

func (x *IConn) options() (opt *options.ClientOptions) {
	opt = options.Client().
		ApplyURI(fmt.Sprintf("mongodb://%s", x.Host)).
		SetReadConcern(readconcern.Majority()).
		SetWriteConcern(writeconcern.Majority()).
		SetAuth(options.Credential{
			Username: x.Username,
			Password: x.Password,
		}).
		SetBSONOptions(&options.BSONOptions{
			UseLocalTimeZone: false,
		})

	opt.Registry = bson.DefaultRegistry

	if x.SetRegistry != nil {
		opt.Registry = x.SetRegistry(opt.Registry)
	}
	return
}
