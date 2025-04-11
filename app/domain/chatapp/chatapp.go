package chatapp

import (
	"context"

	"github.com/gradientsearch/gus/app/api/errs"
)

type App struct {
}

func (a *App) Conversation(ctx context.Context, con Conversation) (Conversation, error) {
	bc, err := toBusConversation(con)
	if err != nil {
		return Conversation{}, errs.New(errs.FailedPrecondition, err)
	}

	ac, err := toAppConversation(bc)
	if err != nil {
		return Conversation{}, errs.New(errs.Internal, err)
	}

	return ac, nil
}
