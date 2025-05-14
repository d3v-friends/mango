package mgConn

import (
	"context"
	"fmt"
	"github.com/d3v-friends/go-tools/fnPointer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
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

func (x *ConnectArgs) Opts() (opt *options.ClientOptions) {
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
		}).
		SetDirect(true)

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

func NewRegistry(codecs ...CodecRegistry) (registry *bsoncodec.Registry) {
	registry = bson.NewRegistry()
	for _, codecRegistry := range codecs {
		registry = codecRegistry(registry)
	}
	return registry
}

func AppendRegistry(
	opt *options.ClientOptions,
	registries ...CodecRegistry,
) *options.ClientOptions {
	if fnPointer.IsNil(opt.Registry) {
		opt.Registry = bson.NewRegistry()
	}

	for _, registry := range registries {
		opt.Registry = registry(opt.Registry)
	}

	return opt
}

func Connect(
	ctx context.Context,
	i *ConnectArgs,
) (client *mongo.Client, err error) {
	if client, err = mongo.Connect(ctx, i.Opts()); err != nil {
		return
	}

	if err = client.Ping(ctx, nil); err != nil {
		return
	}

	return
}
