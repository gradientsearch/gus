// Package all binds all the routes into the specified app.
package all

import (
	"github.com/gradientsearch/gus/api/http/api/mux"
	"github.com/gradientsearch/gus/api/http/domain/checkapi"

	"github.com/gradientsearch/gus/foundation/web"
)

// Routes constructs the add value which provides the implementation of
// of RouteAdder for specifying what routes to bind to this instance.
func Routes() add {
	return add{}
}

type add struct{}

// Add implements the RouterAdder interface.
func (add) Add(app *web.App, cfg mux.Config) {

	checkapi.Routes(app, checkapi.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
	})
}
