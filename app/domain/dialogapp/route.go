package dialogapp

import (
	"net/http"

	"github.com/gradientsearch/gus/app/sdk/auth"
	"github.com/gradientsearch/gus/app/sdk/authclient"
	"github.com/gradientsearch/gus/app/sdk/mid"
	"github.com/gradientsearch/gus/business/domain/dialogbus"
	"github.com/gradientsearch/gus/business/domain/userbus"
	"github.com/gradientsearch/gus/foundation/logger"
	"github.com/gradientsearch/gus/foundation/web"
)

type Config struct {
	Log        *logger.Logger
	UserBus    *userbus.Business
	DialogBus  *dialogbus.Business
	AuthClient *authclient.Client
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	ruleAdminOrSubject := mid.Authorize(cfg.AuthClient, auth.RuleAdminOrSubject)

	api := newApp(*cfg.DialogBus, cfg.Log)

	app.HandlerFunc(http.MethodPost, version, "/dialog", api.create, authen, ruleAdminOrSubject)
}
