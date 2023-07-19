package mango

import (
	"context"
	"fmt"
	"github.com/d3v-friends/mango/migrate"
	"github.com/d3v-friends/mango/models"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type (
	Client struct {
		client   *mongo.Client
		database *mongo.Database
	}
)

func NewClient(opt *ClientOpt) (res *Client, err error) {
	ctx := context.TODO()
	res = &Client{}

	if res.client, err = mongo.Connect(ctx, opt.Options()); err != nil {
		return
	}

	if err = res.client.Ping(ctx, readpref.Primary()); err != nil {
		return
	}

	res.database = res.client.Database(opt.Database)

	return
}

func (x *Client) Database() *mongo.Database {
	return x.database
}

func (x *Client) Client() *mongo.Client {
	return x.client
}

func (x *Client) Migrate(
	ctx context.Context,
	models ...models.IfMigrateModel,
) error {
	return migrate.V1(ctx, x.database, models...)
}

type ClientOpt struct {
	Host     string
	Username string
	Password string
	Database string
	Registry *bsoncodec.Registry
}

func (x *ClientOpt) Options() (opt *options.ClientOptions) {
	opt = options.Client().SetReadConcern(readconcern.Majority()).
		SetWriteConcern(writeconcern.Majority()).
		ApplyURI(fmt.Sprintf("mongodb://%s", x.Host)).
		SetAuth(options.Credential{
			Username: x.Username,
			Password: x.Password,
		})

	if x.Registry != nil {
		opt.SetRegistry(opt.Registry)
	}

	return
}
