package llama

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/gradientsearch/gus/business/domain/chatbus"
)

type StorerMock struct{}

func (s *StorerMock) QueryById(ctx context.Context, userID uuid.UUID, conID uuid.UUID) (chatbus.Conversation, error) {
	return chatbus.Conversation{}, nil
}

func TestQueryLLM(t *testing.T) {
	llama := &Llama{
		BaseURL: "http://localllm.dev:11434",
		Client:  &http.Client{},
		Model:   "llama3.2",
		Stream:  false,
	}

	ms := []chatbus.Message{}

	m1 := chatbus.Message{}
	m1.Role = chatbus.RoleUser
	m1.Content = "My name is stephen!"
	m1.ID = uuid.New()

	ms = append(ms, m1)
	_, err := llama.Chat(ms)
	if err != nil {
		t.Fatalf("FAILED: \tShould be able to query LLM: %s", err)
	}
	t.Logf("SUCCESS: \tShould be able to query LLM")
}
