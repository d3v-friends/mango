package mango

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

/* ------------------------------------------------------------------------------------------------------------ */
// Mango documents

type MigrateData struct {
	ID           primitive.ObjectID `bson:"_id"`
	CollectionNm string             `bson:"collectionNm"`
	Migrated     []*MigrateResult   `bson:"migrated"`
}

type MigrateResult struct {
	Index      int       `bson:"index"`
	Memo       string    `bson:"memo"`
	MigratedAt time.Time `bson:"migratedAt"`
}

func (x *MigrateData) HookTime() []TimeType {
	return []TimeType{
		TimeTypeCreatedAt,
		TimeTypeUpdatedAt,
	}
}

/* ------------------------------------------------------------------------------------------------------------ */
// Migrate

type IfMigrateModel interface {
	IfDocData
	MigrateList() FnMigrateList
}

type FnMigrate func(ctx context.Context, col *mongo.Collection) (memo string, err error)
type FnMigrateList []FnMigrate

func Migrate(
	ctx context.Context,
	mango *Mango,
	models ...IfMigrateModel,
) (err error) {
	panic("not impl")
}

/* ------------------------------------------------------------------------------------------------------------ */
