package cores

import (
	"context"
	"github.com/d3v-friends/mango/opt"
	"go.mongodb.org/mongo-driver/mongo"
)

func Connect(opt *opt.Connect, iCtx ...context.Context) (client *Client, err error) {
	var ctx context.Context
	if 0 < len(iCtx) {
		ctx = iCtx[0]
	} else {
		ctx = context.TODO()
	}

	client = &Client{}
	if client.client, err = mongo.Connect(ctx, opt.Options()); err != nil {
		return
	}

	return
}
