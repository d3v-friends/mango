package mgquery_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/d3v-friends/mango/v2/mgmigrate"
	"github.com/d3v-friends/mango/v2/tester"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestFind(test *testing.T) {
	loadEnv(test)
	var tool = tester.NewTool(test)
	tool.Truncate(test)
	tool.Migrate(test)

	test.Run("findOne", func(t *testing.T) {
		var ctx = tool.Context()

	})

}

func CreateDummy(size int) (ls []*TestModel) {
	ls = make([]*TestModel, size)
	var now = time.Now()
	for i := range ls {
		ls[i] = &TestModel{
			Id:        bson.NewObjectID(),
			Name:      gofakeit.Username(),
			Age:       gofakeit.Int64(),
			CreatedAt: now.Add(time.Hour * -time.Duration(i)),
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

var migrates = mgmigrate.Steps{
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
	Id        bson.ObjectID `bson:"_id"`
	Name      string        `bson:"name"`
	Age       int64         `bson:"age"`
	CreatedAt time.Time     `bson:"createdAt"`
}

func (x TestModel) GetColNm() string {
	return TestModelColNm
}

func (x TestModel) GetMigrates() mgmigrate.Steps {
	return migrates
}

func (x TestModel) IsSame(test *testing.T, v TestModel) {
	assert.Equal(test, x.Id, v.Id)
	assert.Equal(test, x.Name, v.Name)
	assert.Equal(test, x.Age, v.Age)
	assert.Equal(test, x.CreatedAt.Truncate(time.Millisecond).UTC(), v.CreatedAt.Truncate(time.Millisecond).UTC())
}

type TestModelFilter struct {
	Id        *tester.ObjectIdArgs
	Name      *tester.StringArgs
	Age       *tester.Int64Args
	CreatedAt *tester.TimeArgs
}

type TestModelSorter struct {
	Id        *tester.Sorter       `bson:"_id"`
	Name      *tester.Sorter       `bson:"name"`
	Age       *tester.Sorter       `bson:"age"`
	CreatedAt *tester.Sorter       `bson:"createdAt"`
	Data      *TestModelDataSorter `bson:"data"`
}

type TestModelDataSorter struct {
	Title *tester.Sorter `bson:"title"`
}
