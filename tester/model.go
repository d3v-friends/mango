package tester

import (
	"context"
	"testing"
	"time"

	"github.com/d3v-friends/mango/v2/mgmigrate"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
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

func (x Sample) IsSame(test *testing.T, v *Sample) {
	assert.Equal(test, x.Id.Hex(), v.Id.Hex())
	assert.Equal(test, x.Name, v.Name)
	assert.Equal(test, x.Decimal.String(), v.Decimal.String())
	assert.Equal(test, x.CreatedAt.Truncate(time.Millisecond).UTC(), v.CreatedAt.Truncate(time.Millisecond).UTC())
}

/* ---------------------------------------------------------------------- */

type SampleFilter struct {
	Id        *ObjectIdArgs
	Name      *StringArgs
	Decimal   *DecimalArgs
	CreatedAt *TimeArgs
}

type SampleSorter struct {
	Id        *Sorter
	Name      *Sorter
	Decimal   *Sorter
	CreatedAt *Sorter
}
