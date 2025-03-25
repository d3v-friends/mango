package mgConn

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type ConnectArgs struct {
	Host      string
	Username  string
	Password  string
	Codec     []CodecRegistry
	LogOption *options.LoggerOptions
	Monitor   *event.CommandMonitor
}

type CodecRegistry func(*bsoncodec.Registry) *bsoncodec.Registry

func (x *ConnectArgs) host() string {
	return fmt.Sprintf("mongodb://%s", x.Host)
}

func (x *ConnectArgs) opt() (opt *options.ClientOptions) {
	opt = options.Client().
		ApplyURI(x.host()).
		SetReadConcern(readconcern.Majority()).
		SetWriteConcern(writeconcern.Majority()).
		SetAuth(options.Credential{
			Username: x.Username,
			Password: x.Password,
		}).
		SetBSONOptions(&options.BSONOptions{
			UseLocalTimeZone: false,
		})

	if x.LogOption != nil {
		opt.SetLoggerOptions(x.LogOption)
	}

	if x.Monitor != nil {
		opt.SetMonitor(x.Monitor)
	}

	opt.Registry = bson.NewRegistry()

	if len(x.Codec) != 0 {
		for _, registry := range x.Codec {
			opt.Registry = registry(opt.Registry)
		}
	}

	return
}

func Connect(
	ctx context.Context,
	i *ConnectArgs,
) (client *mongo.Client, err error) {
	if client, err = mongo.Connect(ctx, i.opt()); err != nil {
		return
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return
	}

	return
}
