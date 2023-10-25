package mango

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/d3v-friends/go-pure/fnPanic"
	"github.com/d3v-friends/mango/m_migrate"
	"github.com/d3v-friends/mango/m_tx"
	"go.mongodb.org/mongo-driver/mongo"
)

type Mango struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func (x *Mango) Migrate(ctx context.Context, models ...m_migrate.IfMigrateModel) (err error) {
	return m_migrate.Migrate(ctx, x.DB, models...)
}

func (x *Mango) Tx(ctx context.Context, fn m_tx.FnTx) (err error) {
	return m_tx.Transact(ctx, x.DB, fn)
}

/* ------------------------------------------------------------------------------------------------------------ */

const ctxMango = "CTX_MANGO"

func SetMango(ctx context.Context, m *Mango) context.Context {
	return context.WithValue(ctx, ctxMango, m)
}

func GetMango(ctx context.Context) (m *Mango, err error) {
	var isOk bool
	if m, isOk = ctx.Value(ctxMango).(*Mango); !isOk {
		err = fmt.Errorf(
			"not found *mango.Mango in context: context=%s",
			fnPanic.OnValue(json.Marshal(ctx)),
		)
		return
	}
	return
}

func GetMangoP(ctx context.Context) (m *Mango) {
	var err error
	if m, err = GetMango(ctx); err != nil {
		panic(err)
	}
	return
}

/* ------------------------------------------------------------------------------------------------------------ */

func NewMango(
	i *IConn,
	databaseNm string,
	ctxs ...context.Context,
) (res *Mango, err error) {
	res = &Mango{}

	if res.Client, err = NewClient(i, ctxs...); err != nil {
		return
	}

	res.DB = res.Client.Database(databaseNm)

	return
}
