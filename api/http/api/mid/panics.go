package mid

import (
	"context"
	"net/http"

	"github.com/gradientsearch/gus/app/api/mid"
	"github.com/gradientsearch/gus/foundation/web"
)

// Panics executes the panic middleware functionality.
func Panics() web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			hdl := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			return mid.Panics(ctx, hdl)
		}

		return h
	}

	return m
}
