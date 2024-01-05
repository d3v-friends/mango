package v1

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

func NewClient(i *IConn, iCtx ...context.Context) (res *mongo.Client, err error) {
	if i == nil {
		err = fmt.Errorf("IConn is empty value")
		return
	}

	var ctx = context.TODO()
	if len(iCtx) == 1 {
		ctx = iCtx[0]
	}

	if res, err = mongo.Connect(ctx, i.options()); err != nil {
		return
	}

	if err = res.Ping(ctx, readpref.Primary()); err != nil {
		return
	}

	return
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
