package conversationbus

import "github.com/google/uuid"

type Conversation struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

type NewConversation struct {
	ID     uuid.UUID
	UserID uuid.UUID
}
