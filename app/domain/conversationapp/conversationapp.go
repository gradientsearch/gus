package conversationapp

import (
	"context"
	"net/http"

	"github.com/gradientsearch/gus/app/sdk/errs"
	"github.com/gradientsearch/gus/business/domain/conversationbus"
	"github.com/gradientsearch/gus/foundation/logger"
	"github.com/gradientsearch/gus/foundation/web"
)

type App struct {
	conversationBus conversationbus.Business
	log             *logger.Logger
}

func newApp(conversationBus conversationbus.Business, log *logger.Logger) *App {
	return &App{
		conversationBus: conversationBus,
		log:             log,
	}
}

func (a *App) create(ctx context.Context, r *http.Request) web.Encoder {
	newBus, err := toBusNewConversation(ctx)
	if err != nil {
		return errs.New(errs.FailedPrecondition, err)
	}

	c, err := a.conversationBus.Create(ctx, newBus)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	app, err := toAppConversation(c)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return app
}
