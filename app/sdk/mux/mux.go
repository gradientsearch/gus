// Package mux provides support to bind domain level routes
// to the application mux.
package mux

import (
	"embed"
	"net/http"

	"github.com/gradientsearch/gus/app/sdk/auth"
	"github.com/gradientsearch/gus/app/sdk/authclient"
	"github.com/gradientsearch/gus/app/sdk/mid"
	"github.com/gradientsearch/gus/business/domain/conversationbus"
	"github.com/gradientsearch/gus/business/domain/dialogbus"
	"github.com/gradientsearch/gus/business/domain/userbus"
	"github.com/gradientsearch/gus/foundation/logger"
	"github.com/gradientsearch/gus/foundation/web"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
)

// StaticSite represents a static site to run.
type StaticSite struct {
	react      bool
	static     embed.FS
	staticDir  string
	staticPath string
}

// Options represent optional parameters.
type Options struct {
	corsOrigin []string
	sites      []StaticSite
}

// WithCORS provides configuration options for CORS.
func WithCORS(origins []string) func(opts *Options) {
	return func(opts *Options) {
		opts.corsOrigin = origins
	}
}

// WithFileServer provides configuration options for file server.
func WithFileServer(react bool, static embed.FS, dir string, path string) func(opts *Options) {
	return func(opts *Options) {
		opts.sites = append(opts.sites, StaticSite{
			react:      react,
			static:     static,
			staticDir:  dir,
			staticPath: path,
		})
	}
}

// GusConfig contains gus service specific config.
type GusConfig struct {
	AuthClient *authclient.Client
}

// AuthConfig contains auth service specific config.
type AuthConfig struct {
	Auth *auth.Auth
}

type BusConfig struct {
	UserBus         *userbus.Business
	ConversationBus *conversationbus.Business
	DialogBus       *dialogbus.Business
}

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build      string
	Log        *logger.Logger
	DB         *sqlx.DB
	Tracer     trace.Tracer
	BusConfig  BusConfig
	GusConfig  GusConfig
	AuthConfig AuthConfig
}

// RouteAdder defines behavior that sets the routes to bind for an instance
// of the service.
type RouteAdder interface {
	Add(app *web.App, cfg Config)
}

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(cfg Config, routeAdder RouteAdder, options ...func(opts *Options)) http.Handler {
	app := web.NewApp(
		cfg.Log.Info,
		cfg.Tracer,
		mid.Otel(cfg.Tracer),
		mid.Logger(cfg.Log),
		mid.Errors(cfg.Log),
		mid.Metrics(),
		mid.Panics(),
	)

	var opts Options
	for _, option := range options {
		option(&opts)
	}

	if len(opts.corsOrigin) > 0 {
		app.EnableCORS(opts.corsOrigin)
	}

	routeAdder.Add(app, cfg)

	for _, site := range opts.sites {
		if site.react {
			app.FileServerReact(site.static, site.staticDir, site.staticPath)
		} else {
			app.FileServer(site.static, site.staticDir, site.staticPath)
		}
	}

	return app
}
