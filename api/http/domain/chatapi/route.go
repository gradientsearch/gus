package chatapi

import (
	"github.com/gradientsearch/gus/app/domain/chatapp"
	"github.com/gradientsearch/gus/foundation/web"
)

// Routes adds specific routes for this group.
func Routes(app *web.App, chat *chatapp.App) {
	api := newAPI(chat)

	app.HandleFunc("POST /conversation", api.conversation)

}
