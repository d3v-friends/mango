package tester

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/d3v-friends/go-tools/fnEnv"
	"github.com/d3v-friends/go-tools/fnLogger"
	"github.com/d3v-friends/mango/v2/mgconn"
	"github.com/d3v-friends/mango/v2/mgctx"
	"github.com/d3v-friends/mango/v2/mgmigrate"
	"github.com/d3v-friends/mango/v2/mgregistry"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readconcern"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"
)

// Tool
// mango 기능을 사용해서 툴을 만들면 안됀다.
// 기능들을 테스트 하려면 기본 mongo-driver 기능만으로 만들어야 한다.
type Tool struct {
	DB *mongo.Database
}

func NewTool(t *testing.T) (tool *Tool) {
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
		SetMonitor(mgconn.NewMonitor(fnLogger.NewLogger(fnLogger.LogLevelInfo))).
		SetReadConcern(readconcern.Majority()).
		SetWriteConcern(writeconcern.Majority()).
		SetDirect(true)

	opt = mgconn.AppendRegistry(
		opt,
		mgregistry.DecimalRegistry,
	)

	var client, err = mongo.Connect(opt)
	assert.NoError(t, err)

	tool.DB = client.Database(fnEnv.String("MONGO_DATABASE"))

	return
}

func (x *Tool) Context() (ctx context.Context) {
	ctx = context.TODO()
	ctx = fnLogger.SetID(ctx)
	ctx = mgctx.SetDatabase(ctx, x.DB)
	return
}

func (x *Tool) Truncate(t *testing.T) {
	var err = x.DB.Drop(x.Context())
	assert.NoError(t, err)
}

func (x *Tool) Migrate(t *testing.T, models ...mgmigrate.MigratedModel) {
	models = append(models, Sample{})
	assert.NoError(t, mgmigrate.Do(
		x.Context(),
		x.DB,
		models...,
	))
}

func (x *Tool) NewSample() (sample *Sample) {
	sample = &Sample{
		Id:        bson.NewObjectID(),
		Name:      gofakeit.Name(),
		Decimal:   decimal.NewFromInt(rand.Int63()),
		CreatedAt: time.Now(),
	}
	return
}

func (x *Tool) NewSamples(count int) []*Sample {
	var samples = make([]*Sample, count)
	for i := 0; i < count; i++ {
		samples[i] = x.NewSample()
	}
	return samples
}

func (x *Tool) TruncateSampleCollection(t *testing.T) {
	var ctx = x.Context()
	var err = x.DB.Collection(ColNm).Drop(ctx)
	assert.NoError(t, err)
	x.Migrate(t, Sample{})
}
