// Package all binds all the routes into the specified app.
package all

import (
	"github.com/gradientsearch/gus/app/domain/chatapp"
	"github.com/gradientsearch/gus/app/domain/checkapp"
	"github.com/gradientsearch/gus/app/domain/homeapp"
	"github.com/gradientsearch/gus/app/domain/productapp"
	"github.com/gradientsearch/gus/app/domain/rawapp"
	"github.com/gradientsearch/gus/app/domain/tranapp"
	"github.com/gradientsearch/gus/app/domain/userapp"
	"github.com/gradientsearch/gus/app/domain/vproductapp"
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
	checkapp.Routes(app, checkapp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})

	chatapp.Routes(app, chatapp.Config{
		Log:        cfg.Log,
		UserBus:    cfg.BusConfig.UserBus,
		ChatBus:    cfg.BusConfig.ChatBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})

	homeapp.Routes(app, homeapp.Config{
		Log:        cfg.Log,
		HomeBus:    cfg.BusConfig.HomeBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})

	productapp.Routes(app, productapp.Config{
		Log:        cfg.Log,
		ProductBus: cfg.BusConfig.ProductBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})

	rawapp.Routes(app)

	tranapp.Routes(app, tranapp.Config{
		Log:        cfg.Log,
		DB:         cfg.DB,
		UserBus:    cfg.BusConfig.UserBus,
		ProductBus: cfg.BusConfig.ProductBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})

	userapp.Routes(app, userapp.Config{
		Log:        cfg.Log,
		UserBus:    cfg.BusConfig.UserBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})

	vproductapp.Routes(app, vproductapp.Config{
		Log:         cfg.Log,
		UserBus:     cfg.BusConfig.UserBus,
		VProductBus: cfg.BusConfig.VProductBus,
		AuthClient:  cfg.SalesConfig.AuthClient,
	})
}
