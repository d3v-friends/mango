package tester

import (
	"context"
	"fmt"
	"testing"

	"github.com/d3v-friends/go-tools/fnEnv"
	"github.com/d3v-friends/go-tools/fnLogger"
	"github.com/d3v-friends/mango/mgCodec"
	"github.com/d3v-friends/mango/mgConn"
	"github.com/d3v-friends/mango/mgCtx"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type Tool struct {
	DB *mongo.Database
}

func NewTool(t *testing.T) (tool *Tool) {
	var err = fnEnv.Load("../.env")
	assert.NoError(t, err)

	tool = &Tool{}

	var opt = options.Client().
		ApplyURI(fmt.Sprintf(
			"mongodb://%s:%d",
			fnEnv.String("MONGO_HOST"),
			fnEnv.Int("MONGO_PORT"),
		)).
		SetAuth(options.Credential{
			Username: fnEnv.String("MONGO_USERNAME"),
			Password: fnEnv.String("MONGO_PASSWORD"),
		}).
		SetMonitor(mgConn.NewMonitor(fnLogger.NewLogger(fnLogger.LogLevelInfo))).
		SetReadConcern(readconcern.Majority()).
		SetWriteConcern(writeconcern.Majority()).
		SetDirect(true)

	opt = mgConn.AppendRegistry(
		opt,
		mgCodec.DecimalRegistry,
	)

	var client *mongo.Client
	client, err = mongo.Connect(context.TODO(), opt)
	assert.NoError(t, err)

	tool.DB = client.Database(fnEnv.String("MONGO_DATABASE"))

	return
}

func (x *Tool) Context() (ctx context.Context) {
	ctx = context.TODO()
	ctx = fnLogger.SetID(ctx)
	ctx = mgCtx.SetDB(ctx, x.DB)
	return
}
