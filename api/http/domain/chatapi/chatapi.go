// Package chatapi maintains the web based api for system access.
package chatapi

import (
	"context"
	"net/http"
	"time"

	"github.com/gradientsearch/gus/app/api/errs"
	"github.com/gradientsearch/gus/app/domain/chatapp"
	"github.com/gradientsearch/gus/foundation/logger"
	"github.com/gradientsearch/gus/foundation/web"
)

type api struct {
	chatApp *chatapp.App
	log     *logger.Logger
}

func newAPI(chat *chatapp.App, log *logger.Logger) *api {
	return &api{
		chatApp: chat,
		log:     log,
	}
}

func (api *api) conversation(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	var c chatapp.Conversation
	if err := web.Decode(r, &c); err != nil {
		return errs.New(errs.FailedPrecondition, err)
	}

	cr, err := api.chatApp.Conversation(ctx, c)

	if err != nil {
		return errs.New(errs.Internal, err)
	}

	statusCode := http.StatusOK
	return web.Respond(ctx, w, cr, statusCode)
}
