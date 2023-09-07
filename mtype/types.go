package mtype

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	IfModel IfMigrateModel

	IfTrxModel interface {
		GetID() primitive.ObjectID
		SetID(id primitive.ObjectID)
	}

	IfMigrateModel interface {
		IfColNm
		GetMigrateList() FnMigrateList
	}

	FnMigrate     func(ctx context.Context, collection *mongo.Collection) (migrationNm string, err error)
	FnMigrateList []FnMigrate

	IfFilter interface {
		IfColNm
		GetFilter() any
	}

	IfColNm interface {
		GetCollectionNm() string
	}

	IfFilterPager interface {
		IfFilter
		IfPager
	}

	IfPager interface {
		GetPage() int64
		GetSize() int64
	}

	ResultList[MODEL any] struct {
		Page  int64
		Size  int64
		Total int64
		List  []*MODEL
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

type iFilterPager struct {
	Page         int64
	Size         int64
	CollectionNm string
	Filter       any
}

func NewFilterPager(
	collectionNm string,
	filter any,
	page, size int64,
) IfFilterPager {
	return &iFilterPager{
		Page:         page,
		Size:         size,
		CollectionNm: collectionNm,
		Filter:       filter,
	}
}

func (x iFilterPager) GetCollectionNm() string {
	return x.CollectionNm
}

func (x iFilterPager) GetFilter() any {
	return x.Filter
}

func (x iFilterPager) GetPage() int64 {
	return x.Page
}

func (x iFilterPager) GetSize() int64 {
	return x.Size
}
