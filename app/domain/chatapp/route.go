package chatapp

import (
	"net/http"

	"github.com/gradientsearch/gus/app/sdk/auth"
	"github.com/gradientsearch/gus/app/sdk/authclient"
	"github.com/gradientsearch/gus/app/sdk/mid"
	"github.com/gradientsearch/gus/business/domain/chatbus"
	"github.com/gradientsearch/gus/business/domain/userbus"
	"github.com/gradientsearch/gus/foundation/logger"
	"github.com/gradientsearch/gus/foundation/web"
)

type Config struct {
	Log        *logger.Logger
	UserBus    *userbus.Business
	ChatBus    *chatbus.Business
	AuthClient *authclient.Client
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	ruleUserOnly := mid.Authorize(cfg.AuthClient, auth.RuleUserOnly)

	api := newApp(*cfg.ChatBus, cfg.Log)

	app.HandlerFunc(http.MethodPost, "POST /conversation", version, api.conversation, authen, ruleUserOnly)
}
