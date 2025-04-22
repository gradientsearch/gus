package llama

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/gradientsearch/gus/business/domain/dialogbus"
)

type StorerMock struct{}

func (s *StorerMock) QueryById(ctx context.Context, userID uuid.UUID, conID uuid.UUID) (dialogbus.Dialog, error) {
	return dialogbus.Dialog{}, nil
}

func TestQueryLLM(t *testing.T) {
	llama := &Llama{
		BaseURL: "http://localllm.dev:11434",
		Client:  &http.Client{},
		Model:   "llama3.2",
		Stream:  false,
	}

	ms := []dialogbus.Message{}

	sys := dialogbus.Message{}
	sys.Role = dialogbus.RoleSystem
	sys.Content = "You are llm being used for testing purposes. I only want you to respond with the following message: ```I’ve received your message, but I’m only able to acknowledge its receipt. Wishing you a great day ahead!"
	sys.ID = uuid.New()

	m1 := dialogbus.Message{}
	m1.Role = dialogbus.RoleUser
	m1.Content = "My name is stephen!"
	m1.ID = uuid.New()

	ms = append(ms, sys, m1)
	m, err := llama.Chat(ms)
	if err != nil {
		t.Fatalf("FAILED: \tShould be able to query LLM: %s", err)
	}
	t.Logf("SUCCESS: \tShould be able to query LLM")

	if m.Content != "I’ve received your message, but I’m only able to acknowledge its receipt. Wishing you a great day ahead!" {
		t.Fatalf("FAILED: \tShould respond with ACK message but was %s", m.Content)
	}

	t.Logf("SUCCESS: \tShould respond with ACK message")
}
