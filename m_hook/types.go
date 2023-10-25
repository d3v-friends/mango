package m_hook

import (
	"encoding/json"
	"fmt"
	"github.com/d3v-friends/go-pure/fnPanic"
	"time"
)

type (
	// IfCreatedAtModel CreatedAt hook
	IfCreatedAtModel interface {
		GetCreatedAt() time.Time
		SetCreatedAt(t time.Time)
	}

	// IfUpdatedAtModel UpdatedAt hook
	IfUpdatedAtModel interface {
		GetUpdatedAt() time.Time
		SetUpdatedAt(v time.Time)
	}
)

func SetUpdatedAt(model any, t ...time.Time) (err error) {
	var parsed IfUpdatedAtModel
	var isOk bool
	if parsed, isOk = model.(IfUpdatedAtModel); !isOk {
		err = fmt.Errorf(
			"this model not implement mhook.IfUpdatedAtModel interface: model=%s",
			fnPanic.OnValue(json.Marshal(model)),
		)
		return
	}

	if len(t) == 1 {
		parsed.SetUpdatedAt(t[0])
	} else {
		parsed.SetUpdatedAt(time.Now())
	}

	return
}

func SetUpdatedAtWithoutErr(model any, t ...time.Time) {
	_ = SetUpdatedAt(model, t...)
}

/* ------------------------------------------------------------------------------------------------------------ */

func SetCreatedAt(model any, t ...time.Time) (err error) {
	var parsed IfCreatedAtModel
	var isOk bool
	if parsed, isOk = model.(IfCreatedAtModel); !isOk {
		err = fmt.Errorf(
			"this model not implement mhook.IfCreatedAtModel interface: model=%s",
			fnPanic.OnValue(json.Marshal(model)),
		)
		return
	}

	if len(t) == 1 {
		parsed.SetCreatedAt(t[0])
	} else {
		parsed.SetCreatedAt(time.Now())
	}

	return
}

func SetCreatedAtWithoutErr(model any, t ...time.Time) {
	_ = SetCreatedAt(model, t...)
}
