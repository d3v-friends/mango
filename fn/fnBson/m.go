package fnBson

import (
	"fmt"
	"github.com/d3v-friends/go-pure/fnParams"
	"go.mongodb.org/mongo-driver/bson"
)

func ChangeMapToM[T any](kv map[string]T, prefix ...string) (res bson.M) {
	res = make(bson.M)
	p := fnParams.Get(prefix)

	format := ""
	if p == "" {
		format = "%s"
	} else {
		format = fmt.Sprintf("%s.%%s", p)
	}

	for key, value := range kv {
		res[fmt.Sprintf(format, key)] = value
	}
	return
}

func ChangeEmptyToM(list []string, prefix ...string) (res bson.M) {
	res = make(bson.M)

	p := fnParams.Get(prefix)

	format := ""
	if p == "" {
		format = "%s"
	} else {
		format = fmt.Sprintf("%s.%%s", p)
	}

	for _, key := range list {
		res[fmt.Sprintf(format, key)] = ""
	}

	return
}

func MergeM(a bson.M, b bson.M, onConflictErr ...bool) (res bson.M, err error) {
	res = make(bson.M)

	onErr := fnParams.Get(onConflictErr)

	for key, value := range a {
		res[key] = value
	}

	for key, value := range b {
		if _, has := res[key]; onErr && has {
			return nil, fmt.Errorf("duplicatedKey: key=%s", key)
		}
		res[key] = value
	}

	return
}
