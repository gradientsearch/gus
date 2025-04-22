package dialogapp

import (
	"context"
	"net/http"

	"github.com/gradientsearch/gus/app/sdk/errs"
	"github.com/gradientsearch/gus/business/domain/dialogbus"
	"github.com/gradientsearch/gus/foundation/logger"
	"github.com/gradientsearch/gus/foundation/web"
)

type App struct {
	dialogbus dialogbus.Business
	log       *logger.Logger
}

func newApp(dialogbus dialogbus.Business, log *logger.Logger) *App {
	return &App{
		dialogbus: dialogbus,
		log:       log,
	}
}

func (a *App) create(ctx context.Context, r *http.Request) web.Encoder {
	var app Dialog
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	bc, err := toBusConversation(ctx, app)
	if err != nil {
		return errs.New(errs.FailedPrecondition, err)
	}

	c, err := a.dialogbus.Create(ctx, bc)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	ac, err := toAppConversation(c)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return ac
}
