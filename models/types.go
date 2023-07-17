package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
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
		SetID(id primitive.ObjectID) primitive.ObjectID
	}

	IfFilter[MODEL any] interface {
		CollectionNm() string
		Filter() *bson.M
		NewModel() *MODEL
	}

	IfMigrateModel interface {
		CollectionNm() string
		MigrateList() FnMigrateList
	}

	FnMigrate     func(ctx context.Context, collection *mongo.Collection) (migrationNm string, err error)
	FnMigrateList []FnMigrate
)

type filter[T any] struct {
	colNm  string
	filter *bson.M
}

func (x *filter[T]) NewModel() *T {
	return new(T)
}

func (x *filter[T]) CollectionNm() string {
	return x.colNm
}

func (x *filter[T]) Filter() *bson.M {
	return x.filter
}

func NewFilter[T any](colNm string, iFilter *bson.M) IfFilter[T] {
	return &filter[T]{
		colNm:  colNm,
		filter: iFilter,
	}
}
