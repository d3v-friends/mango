package tester

import (
	"context"
	"time"

	"github.com/d3v-friends/mango/v2/mgmigrate"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Sample struct {
	Id        bson.ObjectID   `bson:"_id"`
	Name      string          `bson:"name"`
	Decimal   decimal.Decimal `bson:"decimal"`
	CreatedAt time.Time       `bson:"createdAt"`
}

const (
	ColNm          = "samples"
	FieldId        = "_id"
	FieldName      = "name"
	FieldDecimal   = "decimal"
	FieldCreatedAt = "createdAt"
)

var migrates = mgmigrate.Steps{
	func(ctx context.Context, col *mongo.Collection) (memo string, err error) {
		memo = "init indexing"
		_, err = col.Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys: bson.M{
				FieldCreatedAt: -1,
			},
		})
		return
	},
}

func (x Sample) GetColNm() string {
	return ColNm
}

func (x Sample) GetMigrates() mgmigrate.Steps {
	return migrates
}
