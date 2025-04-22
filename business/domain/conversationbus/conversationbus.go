package conversationbus

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
	Content: "You are a llm being used for testing purposes. I only want you to respond with the following message: ```I’ve received your message, but I’m only able to acknowledge its receipt. Wishing you a great day ahead!",
	Order:   0,
}

type Storer interface {
	QueryById(ctx context.Context, id uuid.UUID, conID uuid.UUID) (Conversation, error)
	Create(ctx context.Context, c NewConversation) error
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

// Create creates a new conversation
func (b *Business) Create(ctx context.Context, newBus NewConversation) (NewConversation, error) {
	if err := b.storer.Create(ctx, newBus); err != nil {
		return NewConversation{}, fmt.Errorf("error creating conversation: %w", err)
	}
	return newBus, nil
}
