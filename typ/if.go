package typ

import "go.mongodb.org/mongo-driver/bson"

type (
	Model interface {
		GetColNm() string
		GetMigrate() []FnMigrate
	}

	Pager interface {
		GetSize() int64
		GetPage() int64
	}

	Filter interface {
		GetFilter() (filter bson.M, err error)
	}

	Sorter interface {
		GetSorter() (filter bson.M, err error)
	}

	Query interface {
		GetQuery() (res bson.M, err error)
	}
)
