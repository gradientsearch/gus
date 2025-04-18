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
	Content: "You are a llm being used for testing purposes. I only want you to respond with the following message: ```I’ve received your message, but I’m only able to acknowledge its receipt. Wishing you a great day ahead!",
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
func (b *Business) Conversation(ctx context.Context, usrConvo Conversation) (Conversation, error) {
	hydrateConvo, err := b.hydrate(ctx, usrConvo)
	if err != nil {
		return Conversation{}, err
	}

	llmMessage, err := b.llm.Chat(hydrateConvo.Messages)

	if err != nil {
		return Conversation{}, fmt.Errorf("error querying llm: %w", err)
	}

	err = b.update(ctx, hydrateConvo, usrConvo, llmMessage)
	if err != nil {
		return Conversation{}, fmt.Errorf("error updating conversation: %w", err)
	}

	hydrateConvo.Messages = []Message{llmMessage}
	return hydrateConvo, nil
}

// hydrate creates a new conversation with a system prompt if it's the start of a conversation,
// or returns the existing conversation with the new messages appended.
func (b *Business) hydrate(ctx context.Context, userConvo Conversation) (Conversation, error) {
	var (
		c   Conversation
		err error
	)

	if userConvo.ID.String() == ROOT_CONVERSATION_ID {
		c.ID = uuid.New()
		c.UserID = userConvo.UserID
		c.Messages = []Message{SYSTEM_PROMPT}
		if err := b.storer.Create(ctx, c); err != nil {
			return Conversation{}, fmt.Errorf("error creating conversation: %w", err)
		}
	} else {
		c, err = b.storer.QueryById(ctx, userConvo.UserID, userConvo.ID)
		if err != nil {
			return Conversation{}, fmt.Errorf("error querying conversation: %w", err)
		}
	}

	c.Messages = append(c.Messages, userConvo.Messages...)

	return c, err
}

// Update writes only the new messages to the database and adds the order to the new messages
func (b *Business) update(ctx context.Context, hydratedConvo Conversation, usrConvo Conversation, llmMessage Message) error {
	nextOrder := len(hydratedConvo.Messages) - len(usrConvo.Messages)

	updateConvo := Conversation{}
	updateConvo.ID = hydratedConvo.ID
	updateConvo.ParentMessageID = hydratedConvo.ParentMessageID
	updateConvo.UserID = hydratedConvo.UserID
	updateConvo.Messages = []Message{}
	updateConvo.Messages = append(updateConvo.Messages, usrConvo.Messages...)
	updateConvo.Messages = append(updateConvo.Messages, llmMessage)

	for i := range updateConvo.Messages {
		updateConvo.Messages[i].Order = nextOrder + i
	}

	b.log.Info(ctx, "messages to update", "messages", fmt.Sprintf("%+v", updateConvo.Messages))
	if err := b.storer.Update(ctx, updateConvo); err != nil {
		return fmt.Errorf("error updating conversation: %w", err)
	}

	return nil
}
