package dialogbus

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/gradientsearch/gus/business/sdk/delegate"
	"github.com/gradientsearch/gus/business/sdk/sqldb"
	"github.com/gradientsearch/gus/foundation/logger"
)

var (
	ErrNotFound = errors.New("dialog not found")
	ErrQuery    = errors.New("querying storage")
)

var SYSTEM_PROMPT = Message{
	ID:      uuid.MustParse("00000000-0000-0000-0000-000000000000"),
	Role:    RoleSystem,
	Content: "You are a llm being used for testing purposes. I only want you to respond with the following message: ```I’ve received your message, but I’m only able to acknowledge its receipt. Wishing you a great day ahead!",
	Order:   0,
}

type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)

	QueryById(ctx context.Context, id uuid.UUID, conID uuid.UUID) (Dialog, error)
	Create(ctx context.Context, c Dialog) error
}

type LLM interface {
	SendCompletionRequest(messages []Message) (Message, error)
}

type Business struct {
	log      *logger.Logger
	storer   Storer
	delegate *delegate.Delegate
	llm      LLM
}

// NewBusiness constructs a user business API for use.
func NewBusiness(log *logger.Logger, delegate *delegate.Delegate, storer Storer, llm LLM) *Business {

	return &Business{
		log:      log,
		storer:   storer,
		delegate: delegate,
		llm:      llm,
	}
}

// NewWithTx constructs a new business value that will use the
// specified transaction in any store related calls.
func (b *Business) NewWithTx(tx sqldb.CommitRollbacker) (*Business, error) {
	storer, err := b.storer.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	bus := Business{
		log:      b.log,
		delegate: b.delegate,
		llm:      b.llm,

		storer: storer,
	}

	return &bus, nil
}

// Create hydrates the conversation with existing messages and updates it with
// new user messages and the LLM response.
func (b *Business) Create(ctx context.Context, bus Dialog) (Dialog, error) {
	promptMessages, err := b.preparePromptMessages(ctx, bus)
	if err != nil {
		return Dialog{}, err
	}

	llmResponse, err := b.llm.SendCompletionRequest(promptMessages.Messages)

	if err != nil {
		return Dialog{}, fmt.Errorf("error querying llm: %w", err)
	}

	err = b.create(ctx, bus, promptMessages, llmResponse)
	if err != nil {
		return Dialog{}, fmt.Errorf("error updating dialog: %w", err)
	}

	promptMessages.Messages = []Message{llmResponse}
	return promptMessages, nil
}

// preparePromptMessages
func (b *Business) preparePromptMessages(ctx context.Context, bus Dialog) (Dialog, error) {
	var (
		d   Dialog
		err error
	)

	d, err = b.storer.QueryById(ctx, bus.UserID, bus.ConversationID)
	if err != nil {
		return Dialog{}, fmt.Errorf("error querying dialog: %w", err)
	}

	d.Messages = append(d.Messages, bus.Messages...)

	d.ConversationID = bus.ConversationID
	d.UserID = bus.UserID

	return d, err
}

// create
func (b *Business) create(ctx context.Context, bus Dialog, promptMessages Dialog, llmResponse Message) error {
	nextOrder := len(promptMessages.Messages) - len(bus.Messages)

	newBus := Dialog{
		ConversationID:  promptMessages.ConversationID,
		ParentMessageID: promptMessages.ParentMessageID,
		UserID:          promptMessages.UserID,
		Messages:        append(bus.Messages, llmResponse),
	}

	for i := range newBus.Messages {
		newBus.Messages[i].Order = nextOrder + i
	}

	if err := b.storer.Create(ctx, newBus); err != nil {
		return fmt.Errorf("error creating dialog messages: %w", err)
	}

	return nil
}
