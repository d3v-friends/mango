package opt

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type (
	Connect struct {
		Host     string
		Username string
		Password string
		Registry *bsoncodec.Registry
	}
)

func (x *Connect) Options() (opt *options.ClientOptions) {
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
