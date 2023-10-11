package fnBson

import (
	"fmt"
	"github.com/d3v-friends/go-pure/fnParams"
	"go.mongodb.org/mongo-driver/bson"
)

func ChangeMapToD[T any](kv map[string]T, prefix ...string) (res bson.D) {
	res = make(bson.D, 0)
	p := fnParams.Get(prefix)

	format := ""
	if p == "" {
		format = "%s"
	} else {
		format = fmt.Sprintf("%s.%%s", p)
	}

	for key, value := range kv {
		res = append(res, bson.E{
			Key:   fmt.Sprintf(format, key),
			Value: value,
		})
	}

	return
}

func ChangeEmptyToD(list []string, prefix ...string) (res bson.D) {
	res = make(bson.D, 0)

	p := fnParams.Get(prefix)
	format := ""

	if p == "" {
		format = "%s"
	} else {
		format = fmt.Sprintf("%s.%%s", p)
	}

	for _, key := range list {
		res = append(res, bson.E{
			Key:   fmt.Sprintf(format, key),
			Value: "",
		})

	}

	return
}

func MergeD(a bson.D, b bson.D) (res bson.D, err error) {
	res = make(bson.D, 0)

	for _, elem := range a {
		if hasKeyInD(res, elem.Key) {
			return nil, fmt.Errorf("duplicatedKey: key=%s", elem.Key)
		}
		res = append(res, elem)
	}

	for _, elem := range b {
		if hasKeyInD(res, elem.Key) {
			return nil, fmt.Errorf("duplicatedKey: key=%s", elem.Key)
		}
		res = append(res, elem)
	}

	return
}

func hasKeyInD(a bson.D, key string) bool {
	for _, v := range a {
		if v.Key == key {
			return true
		}
	}
	return false
}
