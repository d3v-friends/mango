package mango

import "go.mongodb.org/mongo-driver/bson"

type (
	IfFilter interface {
		Filter() (bson.M, error)
		ColNm() string
	}

	IfPager interface {
		Page() int64
		Size() int64
	}

	IfUpdate interface {
		Update() (bson.M, error)
	}
)