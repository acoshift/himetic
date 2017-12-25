package app

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/acoshift/header"
	"github.com/acoshift/hime"
	"github.com/acoshift/httprouter"
	"github.com/acoshift/middleware"
	"github.com/acoshift/session"
	"github.com/acoshift/webstatic"
	yaml "gopkg.in/yaml.v2"
)

// Handler creates new app's handler for given config
func Handler(cfg Config) hime.HandlerFactory {
	return func(app hime.App) http.Handler {
		c := &ctrl{
			sessionName: cfg.SessionName,
			db:          cfg.DB,
		}

		// load static
		static := make(map[string]string)
		{
			bs, err := ioutil.ReadFile("static.yaml")
			if err != nil {
				log.Fatalf("app: can not load static.yaml; %v", err)
			}
			err = yaml.Unmarshal(bs, static)
			if err != nil {
				log.Fatalf("app: can not unmarshal static.yaml; %v", err)
			}
		}

		app.
			TemplateFuncs(template.FuncMap{
				"static": func(name string) string {
					fn, ok := static[name]
					if !ok {
						log.Panicf("app: static %s not exists", name)
					}
					return "/-/" + fn
				},
			}).
			Component("_layout.tmpl").
			Template("index", "index.tmpl").
			Minify().
			BeforeRender(c.beforeRender).
			Routes(hime.Routes{
				"index": "/",
			})

		mux := http.NewServeMux()

		router := httprouter.New()
		router.HandleMethodNotAllowed = false
		router.NotFound = hime.Wrap(c.NotFound)

		router.Get(app.Route("index"), hime.Wrap(indexHandler))

		mux.Handle("/", router)
		mux.Handle("/-/", assetsHeaders(http.StripPrefix("/-", webstatic.New("assets"))))
		mux.Handle("/healthz", hime.Wrap(c.Healthz))

		return middleware.Chain(
			corsProtector,
			securityHeaders,
			session.Middleware(session.Config{
				HTTPOnly: true,
				Path:     "/",
				Secure:   session.PreferSecure,
				SameSite: session.SameSiteLax,
				Store:    cfg.SessionStorage,
				Secret:   cfg.SessionSecret,
			}),
		)(mux)
	}
}

func assetsHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(header.CacheControl, "public, max-age=31536000")
		h.ServeHTTP(w, r)
	})
}

func securityHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(header.XFrameOptions, "deny")
		w.Header().Set(header.XXSSProtection, "1; mode=block")
		w.Header().Set(header.XContentTypeOptions, "nosniff")
		h.ServeHTTP(w, r)
	})
}

func corsProtector(h http.Handler) http.Handler {
	return hime.Wrap(func(ctx hime.Context) hime.Result {
		if ctx.Method() == http.MethodOptions {
			return ctx.Status(http.StatusForbidden).StatusText()
		}
		return ctx.Handle(h)
	})
}
