package fnMango

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

type FnSetRegistry func(registry *bsoncodec.Registry) *bsoncodec.Registry

type ConnectArgs struct {
	Host        string
	Username    string
	Password    string
	SetRegistry []FnSetRegistry
}

func (x *ConnectArgs) options() (opt *options.ClientOptions) {
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

	if len(x.SetRegistry) != 0 {
		for _, registry := range x.SetRegistry {
			opt.Registry = registry(opt.Registry)
		}
	}

	return
}

func Connect(i *ConnectArgs, ctxs ...context.Context) (client *mongo.Client, err error) {
	if i == nil {
		err = fmt.Errorf("IConn is empty value")
		return
	}

	var ctx = context.TODO()
	if 0 < len(ctxs) {
		ctx = ctxs[0]
	}

	if client, err = mongo.Connect(ctx, i.options()); err != nil {
		return
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return
	}

	return
}
