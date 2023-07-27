package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	IfModel interface {
		IfMigrateModel
		IfTrxModel
	}

	IfTrxModel interface {
		GetID() primitive.ObjectID
		SetID(id primitive.ObjectID)
	}

	IfMigrateModel interface {
		GetCollectionNm() string
		GetMigrateList() FnMigrateList
	}

	FnMigrate     func(ctx context.Context, collection *mongo.Collection) (migrationNm string, err error)
	FnMigrateList []FnMigrate

	IfFilter interface {
		GetCollectionNm() string
		GetFilter() any
	}
)

type iFilter struct {
	CollectionNm string
	Filter       any
}

func NewFilter(
	collectionNm string,
	filter any,
) IfFilter {
	return &iFilter{
		CollectionNm: collectionNm,
		Filter:       filter,
	}
}

func (x iFilter) GetCollectionNm() string {
	return x.CollectionNm
}

func (x iFilter) GetFilter() any {
	return x.Filter
}
