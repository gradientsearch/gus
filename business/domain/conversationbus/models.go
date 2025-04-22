package conversationbus

import "github.com/google/uuid"

type Conversation struct {
	ID              uuid.UUID `json:"conversationID"`
	Messages        []Message `json:"messages"`
	ParentMessageID uuid.UUID `json:"parentMessageID"`
	UserID          uuid.UUID `json:"userID"`
}

type Message struct {
	ID      uuid.UUID `json:"id"`
	Role    Role      `json:"role"`
	Content string    `json:"content"`
	Order   int       `json:"order"`
}
