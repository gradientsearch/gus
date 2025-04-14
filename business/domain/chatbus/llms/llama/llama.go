package llama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gradientsearch/gus/business/domain/chatbus"
)

// Llama wraps the base URL and provides a Chat method.
type Llama struct {
	BaseURL string
	Client  *http.Client
	Model   string
	Stream  bool
}

// Chat sends a chat request to the Llama model.
func (l *Llama) Chat(messages []chatbus.Message) (chatbus.Message, error) {
	ms := busToLlmMessages(messages)

	cr := ChatRequest{
		Model:    l.Model,
		Stream:   l.Stream,
		Messages: ms,
	}

	jsonData, err := json.Marshal(cr)
	if err != nil {
		return chatbus.Message{}, fmt.Errorf("error marshalling JSON: %w", err)
	}

	req, err := http.NewRequest("POST", l.BaseURL+"/api/chat", bytes.NewBuffer(jsonData))
	if err != nil {
		return chatbus.Message{}, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := l.Client.Do(req)
	if err != nil {
		return chatbus.Message{}, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	var m chatbus.Message
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return chatbus.Message{}, fmt.Errorf("error decoding body: %w", err)
	}

	return m, nil
}
