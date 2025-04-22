package dialogbus

import "github.com/google/uuid"

type Dialog struct {
	ConversationID  uuid.UUID
	Messages        []Message
	ParentMessageID uuid.UUID
	UserID          uuid.UUID
}

type Message struct {
	ID      uuid.UUID
	Role    Role
	Content string
	Order   int
}
