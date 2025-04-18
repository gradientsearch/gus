package chatapi

import (
	"github.com/gradientsearch/gus/api/http/api/mid"
	"github.com/gradientsearch/gus/app/api/auth"
	"github.com/gradientsearch/gus/app/api/authclient"
	"github.com/gradientsearch/gus/app/domain/chatapp"
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

	authen := mid.Authenticate(cfg.Log, cfg.AuthClient)
	ruleAuthorizeUser := mid.AuthorizeUser(cfg.Log, cfg.AuthClient, cfg.UserBus, auth.RuleAdminOrSubject)

	chatApp := chatapp.New(*cfg.ChatBus, cfg.Log)
	api := newAPI(chatApp, cfg.Log)

	app.HandleFunc("POST /conversation", api.conversation, authen, ruleAuthorizeUser)
}
