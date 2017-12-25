package app

import (
	"github.com/acoshift/hime"
)

func indexHandler(ctx hime.Context) hime.Result {
	return ctx.View("index", nil)
}
