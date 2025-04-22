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

type Storer interface {
	QueryById(ctx context.Context, id uuid.UUID, conID uuid.UUID) (Conversation, error)
	Create(ctx context.Context, c NewConversation) error
}

type Business struct {
	log    *logger.Logger
	storer Storer
}

// NewBusiness constructs a user business API for use.
func NewBusiness(log *logger.Logger, storer Storer) *Business {
	return &Business{
		log:    log,
		storer: storer,
	}
}

// Create creates a new conversation
func (b *Business) Create(ctx context.Context, newBus NewConversation) (NewConversation, error) {
	if err := b.storer.Create(ctx, newBus); err != nil {
		return NewConversation{}, fmt.Errorf("error creating conversation: %w", err)
	}
	return newBus, nil
}
