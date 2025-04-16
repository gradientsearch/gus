package chatbus

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/gradientsearch/gus/foundation/logger"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound = errors.New("conversation not found")
)

const ROOT_CONVERSATION_ID = "00000000-0000-0000-0000-000000000000"

var SYSTEM_PROMPT = Message{
	ID:      uuid.MustParse("00000000-0000-0000-0000-000000000000"),
	Role:    RoleSystem,
	Content: "You are llm being used for testing purposes. I only want you to respond with the following message: ```I’ve received your message, but I’m only able to acknowledge its receipt. Wishing you a great day ahead!",
	Order:   0,
}

type Storer interface {
	QueryById(ctx context.Context, id uuid.UUID, conID uuid.UUID) (Conversation, error)
	Create(ctx context.Context, c Conversation) error
	Update(ctx context.Context, c Conversation) error
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
		c.UserID = con.UserID
		c.Messages = []Message{SYSTEM_PROMPT}
		b.log.Info(ctx, "creating conversation", "conversation", fmt.Sprintf("%+v", c))
		if err := b.storer.Create(ctx, c); err != nil {
			return Conversation{}, fmt.Errorf("error creating conversation: %w", err)
		}
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
	c.Messages = []Message{}
	c.Messages = append(c.Messages, con.Messages...)
	c.Messages = append(c.Messages, llmMessage)

	if err := b.storer.Update(ctx, c); err != nil {
		return Conversation{}, fmt.Errorf("error updating conversation: %w", err)
	}

	c.Messages = []Message{llmMessage}
	return c, nil
}
