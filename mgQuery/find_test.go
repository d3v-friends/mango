package mgQuery_test

import (
	"context"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/d3v-friends/go-tools/fnEnv"
	"github.com/d3v-friends/go-tools/fnLogger"
	"github.com/d3v-friends/go-tools/fnPanic"
	"github.com/d3v-friends/go-tools/fnPointer"
	"github.com/d3v-friends/go-tools/fnSlice"
	"github.com/d3v-friends/mango/mgConn"
	"github.com/d3v-friends/mango/mgCtx"
	"github.com/d3v-friends/mango/mgMigrate"
	"github.com/d3v-friends/mango/mgQuery"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"testing"
	"time"
)

func TestFind(test *testing.T) {
	assert.NoError(test, fnEnv.Load("../.env"))
	var client = fnPanic.Value(mgConn.Connect(context.TODO(), &mgConn.ConnectArgs{
		Host:     fnEnv.String("MONGO_HOST"),
		Username: fnEnv.String("MONGO_USERNAME"),
		Password: fnEnv.String("MONGO_PASSWORD"),
		Monitor:  mgConn.NewMonitor(fnLogger.NewLogger(fnLogger.LogLevelTrace)),
	}))

	var db = client.Database(fnEnv.String("MONGO_DATABASE"))
	assert.NoError(test, db.Drop(context.TODO()))

	var ctx = mgCtx.SetDB(context.TODO(), db)
	assert.NoError(test, mgMigrate.Do(ctx, db, TestModel{}))

	test.Run("sorter", func(t *testing.T) {
		var ls = CreateDummy(10)
		assert.NoError(t, mgQuery.InsertMany(ctx, ls))

		var loaded, err = mgQuery.Find[TestModel](
			ctx,
			bson.M{},
			[]TestModelSorter{
				{
					Age: fnPointer.Make(SorterASC),
				},
				{
					Name: fnPointer.Make(SorterDESC),
				},
			},
			nil)
		assert.NoError(t, err)

		var sorted = fnSlice.Sort(ls, func(a, b *TestModel) bool {
			return a.Name < b.Name
		})

		sorted = fnSlice.Sort(ls, func(a, b *TestModel) bool {
			return a.Age > b.Age
		})

		for i := 0; i < len(ls); i++ {
			loaded[i].IsSame(t, *sorted[i])
		}
	})
}

func CreateDummy(size int) (ls []*TestModel) {
	ls = make([]*TestModel, size)
	var now = time.Now()
	for i := range ls {
		ls[i] = &TestModel{
			Id:        primitive.NewObjectID(),
			Name:      gofakeit.Username(),
			Age:       gofakeit.IntN(80),
			CreatedAt: now.Add(time.Hour * -time.Duration(i)),
			Data: &TestModelData{
				Title: gofakeit.JobTitle(),
			},
		}
	}
	return
}

const (
	TestModelColNm = "testModels"
	FieldId        = "_id"
	FieldName      = "name"
	FieldAge       = "age"
	FieldCreatedAt = "createdAt"
)

var migrates = mgMigrate.Steps{
	func(ctx context.Context, col *mongo.Collection) (memo string, err error) {
		memo = "init indexing"
		_, err = col.Indexes().CreateMany(ctx, []mongo.IndexModel{
			{
				Keys: bson.D{
					{Key: FieldName, Value: 1},
				},
			},
			{
				Keys: bson.D{
					{Key: FieldAge, Value: 1},
				},
			},
			{
				Keys: bson.D{
					{Key: FieldCreatedAt, Value: -1},
				},
			},
		})
		return
	},
}

type TestModel struct {
	Id        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Age       int                `bson:"age"`
	CreatedAt time.Time          `bson:"createdAt"`
	Data      *TestModelData     `bson:"data"`
}

type TestModelData struct {
	Title string `bson:"title"`
}

func (x TestModel) GetColNm() string {
	return TestModelColNm
}

func (x TestModel) GetMigrates() mgMigrate.Steps {
	return migrates
}

func (x TestModel) IsSame(test *testing.T, v TestModel) {
	assert.Equal(test, x.Id, v.Id)
	assert.Equal(test, x.Name, v.Name)
	assert.Equal(test, x.Age, v.Age)
	assert.Equal(test, x.CreatedAt.Truncate(time.Millisecond).UTC(), v.CreatedAt.Truncate(time.Millisecond).UTC())
}

type Sorter string

const (
	SorterUnknown Sorter = "SU_UNKNOWN"
	SorterASC     Sorter = "SU_ASC"
	SorterDESC    Sorter = "SU_DESC"
)

var SorterAll = []Sorter{
	SorterASC,
	SorterDESC,
}

func NewSorter(str string) (res Sorter) {
	res = Sorter(strings.ToUpper(str))
	if res.IsValid() {
		return
	}
	return SorterUnknown
}

func (x Sorter) IsValid() bool {
	for _, v := range SorterAll {
		if v == x {
			return true
		}
	}
	return false
}

func (x Sorter) String() string {
	return string(x)
}

func (x Sorter) GetDirection() int32 {
	if x == SorterASC {
		return 1
	}
	return -1
}

type TestModelSorter struct {
	Id        *Sorter              `bson:"_id"`
	Name      *Sorter              `bson:"name"`
	Age       *Sorter              `bson:"age"`
	CreatedAt *Sorter              `bson:"createdAt"`
	Data      *TestModelDataSorter `bson:"data"`
}

type TestModelDataSorter struct {
	Title *Sorter `bson:"title"`
}
