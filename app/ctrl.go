package app

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/acoshift/header"
	"github.com/acoshift/hime"
	"github.com/acoshift/session"
)

// ctrl is the app's controller
type ctrl struct {
	sessionName string
	db          *sql.DB
}

func (c *ctrl) GetSession(ctx context.Context) *session.Session {
	return session.Get(ctx, c.sessionName)
}

func (*ctrl) NotFound(ctx hime.Context) hime.Result {
	return ctx.NotFound()
}

func (*ctrl) Healthz(ctx hime.Context) hime.Result {
	return ctx.String("ok")
}

func (c *ctrl) beforeRender(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.GetSession(r.Context()).Flash().Clear()
		w.Header().Set(header.CacheControl, "no-cache, no-store, must-revalidate")
		h.ServeHTTP(w, r)
	})
}
