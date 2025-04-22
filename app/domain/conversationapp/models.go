package conversationapp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/gradientsearch/gus/app/sdk/errs"
	"github.com/gradientsearch/gus/app/sdk/mid"
	"github.com/gradientsearch/gus/business/domain/conversationbus"
)

type Conversation struct {
	ID              string `json:"conversationID"`
	ParentMessageID string `json:"parentMessageID"`
}

// the decoder interface.
func (app *Conversation) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks if the data in the model is considered clean.
func (app Conversation) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

// Encode implements the encoder interface.
func (app Conversation) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

type Message struct {
	ID      string `json:"id"`
	Role    string `json:"role"`
	Content string `json:"content"`
}

func toAppConversation(bus conversationbus.Conversation) (Conversation, error) {
	var app Conversation

	app.ID = bus.ID.String()
	app.ParentMessageID = bus.ParentMessageID.String()

	return app, nil
}

func toBusConversation(ctx context.Context, con Conversation) (conversationbus.Conversation, error) {
	var bus conversationbus.Conversation

	if id, err := uuid.Parse(con.ID); err != nil {
		return conversationbus.Conversation{}, fmt.Errorf("bus ID parse: %w", err)
	} else {
		bus.ID = id
	}

	if id, err := uuid.Parse(con.ParentMessageID); err != nil {
		return conversationbus.Conversation{}, fmt.Errorf("bus ParentMessageID parse: %w", err)
	} else {
		bus.ParentMessageID = id
	}

	if userID, err := mid.GetUserID(ctx); err != nil {
		return conversationbus.Conversation{}, fmt.Errorf("bus userID parse: %w", err)
	} else {
		bus.UserID = userID
	}

	return bus, nil
}
