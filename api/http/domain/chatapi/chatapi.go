// Package chatapi maintains the web based api for system access.
package chatapi

import (
	"context"
	"net/http"
	"time"

	"github.com/gradientsearch/gus/app/api/errs"
	"github.com/gradientsearch/gus/app/domain/chatapp"
	"github.com/gradientsearch/gus/foundation/web"
)

type api struct {
	chatApp *chatapp.App
}

func newAPI(chat *chatapp.App) *api {
	return &api{
		chatApp: chat,
	}
}

func (api *api) conversation(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	var c chatapp.Conversation
	if err := web.Decode(r, &c); err != nil {
		return errs.New(errs.FailedPrecondition, err)
	}

	api.chatApp.Conversation(ctx, c)

	status := "ok"
	statusCode := http.StatusOK

	data := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}

	return web.Respond(ctx, w, data, statusCode)
}
