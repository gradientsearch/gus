package chatapp

import (
	"context"

	"github.com/gradientsearch/gus/app/api/errs"
	"github.com/gradientsearch/gus/business/domain/chatbus"
)

type App struct {
	chatBus chatbus.Business
}

func (a *App) Conversation(ctx context.Context, con Conversation) (Conversation, error) {
	bc, err := toBusConversation(ctx, con)
	if err != nil {
		return Conversation{}, errs.New(errs.FailedPrecondition, err)
	}

	c, err := a.chatBus.Conversation(ctx, bc)
	if err != nil {
		return Conversation{}, errs.New(errs.Internal, err)
	}

	ac, err := toAppConversation(c)
	if err != nil {
		return Conversation{}, errs.New(errs.Internal, err)
	}

	return ac, nil
}
