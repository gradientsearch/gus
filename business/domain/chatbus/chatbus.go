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

// Conversation hydrates the conversation with existing messages and updates it with
// new user messages and the LLM response.
func (b *Business) Conversation(ctx context.Context, con Conversation) (Conversation, error) {
	c, err := b.hydrate(ctx, con)
	if err != nil {
		return Conversation{}, err
	}

	llmMessage, err := b.llm.Chat(c.Messages)

	if err != nil {
		return Conversation{}, fmt.Errorf("error querying llm: %w", err)
	}

	err = b.update(ctx, con, llmMessage)
	if err != nil {
		return Conversation{}, fmt.Errorf("error updating conversation: %w", err)
	}

	c.Messages = []Message{llmMessage}
	return c, nil
}

// hydrate creates a new conversation with a system prompt if it's the start of a conversation,
// or returns the existing conversation with the new messages appended.
func (b *Business) hydrate(ctx context.Context, con Conversation) (Conversation, error) {
	var (
		c   Conversation
		err error
	)

	if con.ID.String() == ROOT_CONVERSATION_ID {
		c.ID = uuid.New()
		c.UserID = con.UserID
		c.Messages = []Message{SYSTEM_PROMPT}
		if err := b.storer.Create(ctx, c); err != nil {
			return Conversation{}, fmt.Errorf("error creating conversation: %w", err)
		}
	} else {
		c, err = b.storer.QueryById(ctx, con.UserID, con.ID)
		if err != nil {
			return Conversation{}, fmt.Errorf("error querying conversation: %w", err)
		}
	}

	c.Messages = append(c.Messages, con.Messages...)

	return c, err
}

func (b *Business) update(ctx context.Context, con Conversation, llmMessage Message) error {
	con.Messages = append(con.Messages, llmMessage)

	if err := b.storer.Update(ctx, con); err != nil {
		return fmt.Errorf("error updating conversation: %w", err)
	}

	return nil
}
