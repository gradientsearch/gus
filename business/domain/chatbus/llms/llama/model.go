package llama

import (
	"fmt"

	"github.com/gradientsearch/gus/business/domain/chatbus"
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
	Model    string  `json:"model"`
	Messages Message `json:"message"`
	Stream   bool    `json:"stream"`
}

func busToLlmMessages(messages []chatbus.Message) []Message {
	cm := make([]Message, len(messages))
	for i := range messages {
		var m Message
		m.Content = messages[i].Content
		m.Role = messages[i].Role.Name()
	}
	return cm
}

func llmToBusMessages(msg Message) (chatbus.Message, error) {
	var m chatbus.Message
	m.Content = msg.Content
	r, err := chatbus.ParseRole(msg.Role)
	if err != nil {
		return chatbus.Message{}, fmt.Errorf("unexpected role: %s", msg.Role)
	}
	m.Role = r
	return m, nil
}
