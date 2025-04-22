package conversationdb

import (
	"github.com/google/uuid"
	"github.com/gradientsearch/gus/business/domain/conversationbus"
)

type conversation struct {
	ConversationID uuid.UUID `db:"conversation_id"`
	UserID         uuid.UUID `db:"user_id"`
}

func toBusConversation(db conversation) (conversationbus.Conversation, error) {
	bus := conversationbus.Conversation{
		ID:     db.ConversationID,
		UserID: db.UserID,
	}
	return bus, nil
}

func toDbConversation(bus conversationbus.NewConversation) conversation {
	db := conversation{
		ConversationID: bus.ID,
		UserID:         bus.UserID,
	}
	return db
}
