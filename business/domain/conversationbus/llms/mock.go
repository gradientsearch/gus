package llms

import (
	"github.com/google/uuid"
	"github.com/gradientsearch/gus/business/domain/conversationbus"
)

// Llama wraps the base URL and provides a Chat method.
type Mock struct {
}

// Chat Mock sends a chat request to the Llama model.
func (m *Mock) Chat(messages []conversationbus.Message) (conversationbus.Message, error) {
	msg := conversationbus.Message{}
	msg.Role = conversationbus.RoleAssistant
	msg.Content = "I’ve received your message, but I’m only able to acknowledge its receipt. Wishing you a great day ahead!"
	msg.ID = uuid.New()

	return msg, nil
}
