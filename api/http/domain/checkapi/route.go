package checkapi

import (
	"github.com/gradientsearch/gus/foundation/logger"
	"github.com/gradientsearch/gus/foundation/web"
	"github.com/jmoiron/sqlx"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
	Log   *logger.Logger
	DB    *sqlx.DB
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	api := newAPI(cfg.Build, cfg.Log)

	app.HandleFuncNoMiddleware("GET /liveness", api.liveness)
	app.HandleFuncNoMiddleware("GET /readiness", api.readiness)
}
