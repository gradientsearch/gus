package dialogdb

import (
	"github.com/google/uuid"
	"github.com/gradientsearch/gus/business/domain/dialogbus"
)

type message struct {
	MessageID      uuid.UUID `db:"message_id"`
	ConversationID uuid.UUID `db:"conversation_id"`
	UserID         uuid.UUID `db:"user_id"`
	Role           string    `db:"role"`
	Content        string    `db:"content"`
	Order          int       `db:"order"`
}

func toBusDialog(db []message) (dialogbus.Dialog, error) {
	bus := dialogbus.Dialog{
		Messages: make([]dialogbus.Message, 0),
	}

	for _, m := range db {
		bm := dialogbus.Message{
			ID:      m.MessageID,
			Content: m.Content,
			Role:    dialogbus.NewRole(m.Role),
			Order:   m.Order,
		}

		bus.Messages = append(bus.Messages, bm)
	}

	return bus, nil
}

func toDbMessages(bus dialogbus.Dialog) []message {
	db := make([]message, len(bus.Messages))
	for i, bm := range bus.Messages {
		m := message{
			MessageID:      bm.ID,
			ConversationID: bus.ConversationID,
			UserID:         bus.UserID,

			Content: bm.Content,
			Role:    bm.Role.Name(),
			Order:   bm.Order,
		}

		db[i] = m
	}

	return db
}
