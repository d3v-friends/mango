package fnMango

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type IConnect struct {
	Host     string
	Username string
	Password string
	Registry *bsoncodec.Registry
}

func (x *IConnect) Options() (opt *options.ClientOptions) {
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

func Connect(ctx context.Context, i *IConnect) (client *mongo.Client, err error) {
	if client, err = mongo.Connect(ctx, i.Options()); err != nil {
		return
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return
	}
	return
}
