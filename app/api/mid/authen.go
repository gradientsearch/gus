package mid

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/gradientsearch/gus/app/api/auth"
	"github.com/gradientsearch/gus/app/api/authclient"
	"github.com/gradientsearch/gus/app/api/errs"
	"github.com/gradientsearch/gus/business/domain/userbus"

	"github.com/gradientsearch/gus/foundation/logger"
)

// Authenticate validates authentication via the auth service.
func Authenticate(ctx context.Context, log *logger.Logger, client *authclient.Client, authorization string, handler Handler) error {
	resp, err := client.Authenticate(ctx, authorization)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	ctx = setUserID(ctx, resp.UserID)
	ctx = setClaims(ctx, resp.Claims)

	return handler(ctx)
}

// Bearer processes JWT authentication logic.
func Bearer(ctx context.Context, ath *auth.Auth, authorization string, handler Handler) error {
	claims, err := ath.Authenticate(ctx, authorization)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	if claims.Subject == "" {
		return errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, no claims")
	}

	subjectID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return errs.New(errs.Unauthenticated, fmt.Errorf("parsing subject: %w", err))
	}

	ctx = setUserID(ctx, subjectID)
	ctx = setClaims(ctx, claims)

	return handler(ctx)
}

// Basic processes basic authentication logic.
func Basic(ctx context.Context, r *http.Request, userBus userbus.Business, handler Handler) error {
	email, pass, ok := parseBasicAuth(r.Header.Get("authorization"))
	if !ok {
		return errs.Newf(errs.Unauthenticated, "invalid Basic auth")
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	usr, err := userBus.Authenticate(ctx, *addr, pass)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usr.ID.String(),
			Issuer:    "gus-system",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: userbus.ParseToString(usr.Roles),
	}

	subjectID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return errs.Newf(errs.Unauthenticated, "parsing subject: %s", err)
	}

	ctx = setUserID(ctx, subjectID)
	ctx = setClaims(ctx, claims)

	return handler(ctx)
}

func parseBasicAuth(auth string) (string, string, bool) {
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Basic" {
		return "", "", false
	}

	c, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", false
	}

	username, password, ok := strings.Cut(string(c), ":")
	if !ok {
		return "", "", false
	}

	return username, password, true
}
