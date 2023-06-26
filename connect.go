package mango

import (
	"context"
	"github.com/d3v-friends/mango/cores"
	"github.com/d3v-friends/mango/opt"
)

func Connect(opt *opt.Connect, iCtx ...context.Context) (*cores.Client, error) {
	return cores.Connect(opt, iCtx...)
}
