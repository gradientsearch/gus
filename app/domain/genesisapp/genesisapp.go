// Package genesisapp maintains the app layer api for the tran domain.
package genesisapp

import (
	"context"
	"net/http"

	"github.com/gradientsearch/gus/app/sdk/errs"
	"github.com/gradientsearch/gus/app/sdk/mid"
	"github.com/gradientsearch/gus/business/domain/conversationbus"
	"github.com/gradientsearch/gus/business/domain/dialogbus"
	"github.com/gradientsearch/gus/foundation/web"
)

type app struct {
	conversationBus *conversationbus.Business
	dialogBus       *dialogbus.Business
}

func newApp(conversationBus *conversationbus.Business, dialogBus *dialogbus.Business) *app {
	return &app{
		conversationBus: conversationBus,
		dialogBus:       dialogBus,
	}
}

// newWithTx constructs a new Handlers value with the domain apis
// using a store transaction that was created via middleware.
func (a *app) newWithTx(ctx context.Context) (*app, error) {
	tx, err := mid.GetTran(ctx)
	if err != nil {
		return nil, err
	}

	conversationBus, err := a.conversationBus.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	dialogBus, err := a.dialogBus.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	app := app{
		conversationBus: conversationBus,
		dialogBus:       dialogBus,
	}

	return &app, nil
}

func (a *app) create(ctx context.Context, r *http.Request) web.Encoder {
	var app NewDialog
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	a, err := a.newWithTx(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	nc, err := toBusNewConversation(ctx)
	if err != nil {
		return errs.New(errs.FailedPrecondition, err)
	}

	c, err := a.conversationBus.Create(ctx, nc)
	if err != nil {
		return errs.New(errs.FailedPrecondition, err)
	}

	nd, err := toBusNewDialog(ctx, app, c.ID)
	if err != nil {
		return errs.New(errs.FailedPrecondition, err)
	}

	bus, err := a.dialogBus.Create(ctx, nd)
	if err != nil {
		return errs.New(errs.FailedPrecondition, err)
	}

	d, err := toAppDialog(bus)
	if err != nil {
		return errs.New(errs.FailedPrecondition, err)
	}

	return d
}
