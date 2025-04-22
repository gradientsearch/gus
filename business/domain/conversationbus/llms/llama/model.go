package llama

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/gradientsearch/gus/business/domain/conversationbus"
)

// Message represents a single chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest is the full request payload.
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

// ChatRequest is the full request payload.
type ChatResponse struct {
	Model   string  `json:"model"`
	Message Message `json:"message"`
	Stream  bool    `json:"stream"`
}

func busToLlmMessages(messages []conversationbus.Message) []Message {
	cm := make([]Message, len(messages))
	for i := range messages {
		var m Message
		m.Content = messages[i].Content
		m.Role = messages[i].Role.Name()
		cm[i] = m
	}
	return cm
}

func llmToBusMessages(msg Message) (conversationbus.Message, error) {
	var m conversationbus.Message
	m.ID = uuid.New()
	m.Content = msg.Content
	r, err := conversationbus.ParseLlmRoles(msg.Role)
	if err != nil {
		return conversationbus.Message{}, fmt.Errorf("unexpected role: %s", msg.Role)
	}
	m.Role = r
	// Will be reordered later. Stub value greater than 0; system prompt uses 0.
	m.Order = 1e5
	return m, nil
}
