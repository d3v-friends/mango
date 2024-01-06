package stDoc

import (
	"context"
	"github.com/d3v-friends/go-pure/fnReflect"
	"github.com/d3v-friends/mango/typ"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type (
	Mango struct {
		Id        primitive.ObjectID `bson:"_id"`
		ColNm     string             `bson:"colNm"`
		NextIdx   int                `bson:"nextIdx"`
		History   []*MangoHistory    `bson:"history"`
		CreatedAt time.Time          `bson:"createdAt"`
		UpdatedAt time.Time          `bson:"updatedAt"`
	}

	MangoHistory struct {
		Memo       string    `bson:"memo"`
		MigratedAt time.Time `bson:"migratedAt"`
	}
)

func (x *Mango) GetColNm() string {
	return "mango"
}

func (x *Mango) GetMigrate() []typ.FnMigrate {
	return []typ.FnMigrate{
		func(ctx context.Context, col *mongo.Collection) (memo string, err error) {
			memo = "init indexing"
			_, err = col.Indexes().CreateMany(ctx, []mongo.IndexModel{
				{
					Keys: bson.D{
						{
							Key:   "colNm",
							Value: 1,
						},
					},
					Options: &options.IndexOptions{
						Unique: fnReflect.ToPointer(true),
					},
				},
			})

			return
		},
	}
}
