package mango

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/d3v-friends/go-pure/fnPanic"
	"github.com/d3v-friends/mango/mMigrate"
	"github.com/d3v-friends/mango/mTx"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Mango struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func (x *Mango) Migrate(ctx context.Context, models ...mMigrate.IfMigrateModel) (err error) {
	return mMigrate.Migrate(ctx, x.DB, models...)
}

func (x *Mango) Tx(ctx context.Context, fn mTx.FnTx) (err error) {
	return mTx.Transact(ctx, x.DB, fn)
}

func (x *Mango) TxWithDelay(ctx context.Context, fn mTx.FnTx, delay time.Duration) (err error) {
	return mTx.TransactWithDelay(ctx, x.DB, fn, delay)
}

func (x *Mango) Truncate(ctx context.Context) error {
	return x.DB.Drop(ctx)
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
	ctxs ...context.Context,
) (res *Mango, err error) {
	res = &Mango{}

	if res.Client, err = NewClient(i, ctxs...); err != nil {
		return
	}

	res.DB = res.Client.Database(i.Database)

	return
}
