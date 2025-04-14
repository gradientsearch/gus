package chatbus

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/gradientsearch/gus/foundation/logger"
)

const ROOT_CONVERSATION_ID = "00000000-0000-0000-0000-000000000000"

type Storer interface {
	QueryById(ctx context.Context, id uuid.UUID, conID uuid.UUID) (Conversation, error)
}

type LLM interface {
	Chat(messages []Message) (Message, error)
}

type Business struct {
	log    *logger.Logger
	storer Storer
	llm    LLM
}

// NewBusiness constructs a user business API for use.
func NewBusiness(log *logger.Logger, storer Storer, llm LLM) *Business {

	return &Business{
		log:    log,
		storer: storer,
		llm:    llm,
	}
}

func (b *Business) Conversation(ctx context.Context, con Conversation) (Conversation, error) {
	var c Conversation
	var err error

	if con.ID.String() == ROOT_CONVERSATION_ID {
		c = Conversation{}
		c.ID = uuid.New()
	} else {
		c, err = b.storer.QueryById(ctx, con.UserID, con.ID)
		if err != nil {
			return Conversation{}, fmt.Errorf("error querying conversation: %w", err)
		}
	}

	// Append new message[s] to existing conversation
	c.Messages = append(c.Messages, con.Messages...)

	llmMessage, err := b.llm.Chat(c.Messages)
	if err != nil {
		return Conversation{}, fmt.Errorf("error querying llm: %w", err)
	}
	c.Messages = append(c.Messages, llmMessage)

	return c, nil
}
