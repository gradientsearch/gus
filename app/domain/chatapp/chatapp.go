package chatapp

import (
	"context"
	"net/http"

	"github.com/gradientsearch/gus/app/sdk/errs"
	"github.com/gradientsearch/gus/business/domain/chatbus"
	"github.com/gradientsearch/gus/foundation/logger"
	"github.com/gradientsearch/gus/foundation/web"
)

type App struct {
	chatBus chatbus.Business
	log     *logger.Logger
}

func newApp(chatBus chatbus.Business, log *logger.Logger) *App {
	return &App{
		chatBus: chatBus,
		log:     log,
	}
}

func (a *App) conversation(ctx context.Context, r *http.Request) web.Encoder {
	var app Conversation
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	bc, err := toBusConversation(ctx, app)
	if err != nil {
		return errs.New(errs.FailedPrecondition, err)
	}

	c, err := a.chatBus.Conversation(ctx, bc)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	ac, err := toAppConversation(c)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return ac
}
