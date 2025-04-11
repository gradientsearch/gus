package chatbus

import "github.com/google/uuid"

type Conversation struct {
	ID              uuid.UUID `json:"conversationID"`
	Messages        []Message `json:"messages"`
	ParentMessageID uuid.UUID `json:"parentMessageID"`
}

type Message struct {
	ID      uuid.UUID `json:"id"`
	Role    Role      `json:"author"`
	Content string    `json:"content"`
}
