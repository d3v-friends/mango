package m_tx

import (
	"context"
	"github.com/d3v-friends/go-pure/fnPanic"
	"github.com/d3v-friends/mango/m_ctx"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type (
	FnTx func(tx *TxDB) (txErr error)
)

func Transact(
	ctx context.Context,
	db *mongo.Database,
	fn FnTx,
) (err error) {
	var txDB = NewTxDB(ctx, db)

	if err = fn(txDB); err == nil {
		fnPanic.On(txDB.commit())
	} else {
		fnPanic.On(txDB.rollback())
	}

	return
}

func IncludeCtx(
	ctx context.Context,
	fn FnTx,
) (err error) {
	var db = m_ctx.GetDBP(ctx)
	return Transact(ctx, db, fn)
}

// TransactWithDelay transact 실험용 function
func TransactWithDelay(
	ctx context.Context,
	db *mongo.Database,
	fn FnTx,
	delay time.Duration,
) (err error) {
	var txDB = &TxDB{
		ctx:    ctx,
		db:     db,
		insert: make([]*insertModel, 0),
		delete: make([]*deleteModel, 0),
		update: make([]*updateModel, 0),
		lock:   make([]*lockModel, 0),
	}

	if err = fn(txDB); err == nil {
		time.Sleep(delay)
		fnPanic.On(txDB.commit())
	} else {
		time.Sleep(delay)
		fnPanic.On(txDB.rollback())
	}

	return
}
