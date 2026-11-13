package mgconn

import (
	"context"
	"fmt"

	"github.com/d3v-friends/go-tools/fnPointer"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readconcern"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"
)

type ConnectArgs struct {
	Host      string
	Username  string
	Password  string
	Codec     []CodecRegistry
	LogOption *options.LoggerOptions
	Monitor   *event.CommandMonitor
}

type CodecRegistry func(*bson.Registry) *bson.Registry

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
		SetDirect(true).
		SetRegistry(bson.NewRegistry())

	if x.LogOption != nil {
		opt.SetLoggerOptions(x.LogOption)
	}

	if x.Monitor != nil {
		opt.SetMonitor(x.Monitor)
	}

	if len(x.Codec) != 0 {
		for _, registry := range x.Codec {
			opt.Registry = registry(opt.Registry)
		}
	}

	return
}

func NewRegistry(codecs ...CodecRegistry) (registry *bson.Registry) {
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
	if client, err = mongo.Connect(i.Opts()); err != nil {
		return
	}

	if err = client.Ping(ctx, nil); err != nil {
		return
	}

	return
}
