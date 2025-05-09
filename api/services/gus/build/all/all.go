// Package all binds all the routes into the specified app.
package all

import (
	"github.com/gradientsearch/gus/app/domain/checkapp"
	"github.com/gradientsearch/gus/app/domain/conversationapp"
	"github.com/gradientsearch/gus/app/domain/dialogapp"
	"github.com/gradientsearch/gus/app/domain/genesisapp"
	"github.com/gradientsearch/gus/app/domain/rawapp"
	"github.com/gradientsearch/gus/app/domain/userapp"
	"github.com/gradientsearch/gus/app/sdk/mux"
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

	// =============================================================================================
	// service specific routes

	checkapp.Routes(app, checkapp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})

	rawapp.Routes(app)

	userapp.Routes(app, userapp.Config{
		Log:        cfg.Log,
		UserBus:    cfg.BusConfig.UserBus,
		AuthClient: cfg.GusConfig.AuthClient,
	})

	// =============================================================================================
	// GUS specific routes

	genesisapp.Routes(app, genesisapp.Config{
		Log:        cfg.Log,
		DB:         cfg.DB,
		AuthClient: cfg.GusConfig.AuthClient,

		DialogBus:       cfg.BusConfig.DialogBus,
		ConversationBus: cfg.BusConfig.ConversationBus,
	})

	conversationapp.Routes(app, conversationapp.Config{
		Log:        cfg.Log,
		AuthClient: cfg.GusConfig.AuthClient,

		UserBus:         cfg.BusConfig.UserBus,
		ConversationBus: cfg.BusConfig.ConversationBus,
	})

	dialogapp.Routes(app, dialogapp.Config{
		Log:        cfg.Log,
		AuthClient: cfg.GusConfig.AuthClient,

		UserBus:   cfg.BusConfig.UserBus,
		DialogBus: cfg.BusConfig.DialogBus,
	})

}
