package conversationdb

import (
	"github.com/google/uuid"
	"github.com/gradientsearch/gus/business/domain/conversationbus"
)

type conversation struct {
	ConversationID  uuid.UUID `db:"conversation_id"`
	ParentMessageID uuid.UUID `db:"parent_message_id"`
	UserID          uuid.UUID `db:"user_id"`
}

func toBusConversation(db conversation) (conversationbus.Conversation, error) {
	bus := conversationbus.Conversation{
		ID:              db.ConversationID,
		ParentMessageID: db.ParentMessageID,
		UserID:          db.UserID,
	}
	return bus, nil
}

func toDbConversation(bus conversationbus.Conversation) conversation {
	db := conversation{
		ConversationID:  bus.ID,
		ParentMessageID: bus.ParentMessageID,
		UserID:          bus.UserID,
	}
	return db
}
