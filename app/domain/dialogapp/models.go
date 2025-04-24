package dialogapp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/gradientsearch/gus/app/sdk/errs"
	"github.com/gradientsearch/gus/app/sdk/mid"
	"github.com/gradientsearch/gus/business/domain/dialogbus"
)

type Dialog struct {
	ConversationID  string    `json:"conversationID"`
	Messages        []Message `json:"messages"`
	ParentMessageID string    `json:"parentMessageID"`
}

// the decoder interface.
func (app *Dialog) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks if the data in the model is considered clean.
func (app Dialog) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

// Encode implements the encoder interface.
func (app Dialog) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

type Message struct {
	ID      string `json:"id"`
	Role    string `json:"role"`
	Content string `json:"content"`
}

func toAppDialog(bus dialogbus.Dialog) (Dialog, error) {
	var app Dialog

	app.ConversationID = bus.ConversationID.String()
	app.ParentMessageID = bus.ParentMessageID.String()

	if m, err := toAppMessages(bus.Messages); err != nil {
		return Dialog{}, err
	} else {
		app.Messages = m
	}

	return app, nil
}

func toAppMessages(bus []dialogbus.Message) ([]Message, error) {
	app := make([]Message, len(bus))
	for i, b := range bus {
		var a Message
		a.ID = b.ID.String()
		a.Role = b.Role.Name()
		a.Content = b.Content
		app[i] = a
	}

	return app, nil
}

func toBusDialog(ctx context.Context, con Dialog) (dialogbus.Dialog, error) {
	var bus dialogbus.Dialog

	if id, err := uuid.Parse(con.ConversationID); err != nil {
		return dialogbus.Dialog{}, fmt.Errorf("bus ID parse: %w", err)
	} else {
		bus.ConversationID = id
	}

	if id, err := uuid.Parse(con.ParentMessageID); err != nil {
		return dialogbus.Dialog{}, fmt.Errorf("bus ParentMessageID parse: %w", err)
	} else {
		bus.ParentMessageID = id
	}

	if mes, err := toBusMessages(con.Messages); err != nil {
		return dialogbus.Dialog{}, err
	} else {
		bus.Messages = mes
	}

	if userID, err := mid.GetUserID(ctx); err != nil {
		return dialogbus.Dialog{}, fmt.Errorf("bus userID parse: %w", err)
	} else {
		bus.UserID = userID
	}

	return bus, nil
}

func toBusMessages(app []Message) ([]dialogbus.Message, error) {
	bus := make([]dialogbus.Message, len(app))

	for i, m := range app {
		var b dialogbus.Message

		if id, err := uuid.Parse(m.ID); err != nil {
			return nil, fmt.Errorf("bus message ID parse: %w", err)
		} else {
			b.ID = id
		}

		if role, err := dialogbus.ParseUserRoles(m.Role); err != nil {
			return nil, fmt.Errorf("bus message Role parse: %w", err)
		} else {
			b.Role = role
		}

		// TODO sanitize content
		b.Content = m.Content

		bus[i] = b
	}

	return bus, nil
}
