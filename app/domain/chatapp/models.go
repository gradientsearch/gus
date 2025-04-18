package chatapp

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/gradientsearch/gus/app/api/mid"
	"github.com/gradientsearch/gus/business/domain/chatbus"
)

type Conversation struct {
	ID              string    `json:"conversationID"`
	Messages        []Message `json:"messages"`
	ParentMessageID string    `json:"parentMessageID"`
}

type Message struct {
	ID      string `json:"id"`
	Role    string `json:"role"`
	Content string `json:"content"`
}

func toAppConversation(bus chatbus.Conversation) (Conversation, error) {
	var app Conversation

	app.ID = bus.ID.String()
	app.ParentMessageID = bus.ParentMessageID.String()

	if m, err := toAppMessages(bus.Messages); err != nil {
		return Conversation{}, err
	} else {
		app.Messages = m
	}

	return app, nil
}

func toAppMessages(bus []chatbus.Message) ([]Message, error) {
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

func toBusConversation(ctx context.Context, con Conversation) (chatbus.Conversation, error) {
	var bus chatbus.Conversation

	if id, err := uuid.Parse(con.ID); err != nil {
		return chatbus.Conversation{}, fmt.Errorf("bus ID parse: %w", err)
	} else {
		bus.ID = id
	}

	if id, err := uuid.Parse(con.ParentMessageID); err != nil {
		return chatbus.Conversation{}, fmt.Errorf("bus ParentMessageID parse: %w", err)
	} else {
		bus.ParentMessageID = id
	}

	if mes, err := toBusMessages(con.Messages); err != nil {
		return chatbus.Conversation{}, err
	} else {
		bus.Messages = mes
	}

	if userID, err := mid.GetUserID(ctx); err != nil {
		return chatbus.Conversation{}, fmt.Errorf("bus userID parse: %w", err)
	} else {
		bus.UserID = userID
	}

	return bus, nil
}

func toBusMessages(app []Message) ([]chatbus.Message, error) {
	bus := make([]chatbus.Message, len(app))

	for i, m := range app {
		var b chatbus.Message

		if id, err := uuid.Parse(m.ID); err != nil {
			return nil, fmt.Errorf("bus message ID parse: %w", err)
		} else {
			b.ID = id
		}

		if role, err := chatbus.ParseUserRoles(m.Role); err != nil {
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
