package fnMigrate

import (
	"github.com/d3v-friends/go-pure/fnEnv"
	"github.com/d3v-friends/go-pure/fnPanic"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestMigrate(test *testing.T) {
	fnPanic.On(fnEnv.ReadFromFile("../env/.env"))

	test.Run("migrate", func(t *testing.T) {

	})
}

type DocTest struct {
	Id primitive.ObjectID `bson:"_id"`
}
