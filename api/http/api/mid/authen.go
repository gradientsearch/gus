package mid

import (
	"context"
	"net/http"

	"github.com/gradientsearch/gus/app/api/auth"
	"github.com/gradientsearch/gus/app/api/authclient"
	"github.com/gradientsearch/gus/app/api/mid"
	"github.com/gradientsearch/gus/foundation/logger"
	"github.com/gradientsearch/gus/foundation/web"
)

// Authenticate validates authentication via the auth service.
func Authenticate(log *logger.Logger, client *authclient.Client) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			hdl := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			return mid.Authenticate(ctx, log, client, r.Header.Get("authorization"), hdl)
		}

		return h
	}

	return m
}

// Bearer processes JWT authentication logic.
func Bearer(ath *auth.Auth) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			hdl := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			return mid.Bearer(ctx, ath, r.Header.Get("authorization"), hdl)
		}

		return h
	}

	return m
}

// Basic processes basic authentication logic.
func Basic(ath *auth.Auth) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			hdl := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			return mid.Basic(ctx, r, *ath.UserBus, hdl)
		}
		return h
	}

	return m
}
