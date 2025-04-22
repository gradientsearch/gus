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
	ID string `json:"conversationID"`
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

func toAppConversation(bus conversationbus.NewConversation) (Conversation, error) {
	app := Conversation{
		ID: bus.ID.String(),
	}
	return app, nil
}

func toBusNewConversation(ctx context.Context) (conversationbus.NewConversation, error) {
	var bus conversationbus.NewConversation

	bus.ID = uuid.New()
	if userID, err := mid.GetUserID(ctx); err != nil {
		return conversationbus.NewConversation{}, fmt.Errorf("bus userID parse: %w", err)
	} else {
		bus.UserID = userID
	}

	return bus, nil
}

// =================================================================================================

type NewConversation struct{}
