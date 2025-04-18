package chatdb

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/gradientsearch/gus/business/domain/chatbus"
)

type conversation struct {
	ConversationID  uuid.UUID `db:"conversation_id"`
	ParentMessageID uuid.UUID `db:"parent_message_id"`
	UserID          uuid.UUID `db:"user_id"`
}

type conversationMessages struct {
	conversation
	message
}

type message struct {
	MessageID uuid.UUID `db:"message_id"`
	Role      string    `db:"role"`
	Content   string    `db:"content"`
	Order     int       `db:"order"`
}

func toBusConversation(dbCon []conversationMessages) (chatbus.Conversation, error) {
	if len(dbCon) < 1 {
		return chatbus.Conversation{}, fmt.Errorf("db: conversation not found error")
	}

	busConvo := chatbus.Conversation{}

	for _, m := range dbCon {
		bm := chatbus.Message{}
		bm.ID = m.MessageID
		bm.Content = m.message.Content
		bm.Role = chatbus.NewRole(m.Role)
		bm.Order = m.Order
		busConvo.Messages = append(busConvo.Messages, bm)
	}

	busConvo.ID = dbCon[0].ConversationID
	busConvo.ParentMessageID = dbCon[0].ParentMessageID
	busConvo.UserID = dbCon[0].UserID

	return busConvo, nil
}

func toDbConversation(busCon chatbus.Conversation) conversation {
	var dbCon = conversation{}
	dbCon.ConversationID = busCon.ID
	dbCon.ParentMessageID = busCon.ParentMessageID
	dbCon.UserID = busCon.UserID
	return dbCon
}

func toDbMessages(busMsgs []chatbus.Message) []message {
	dbMsgs := make([]message, len(busMsgs))
	for i, bm := range busMsgs {
		dm := message{}
		dm.MessageID = bm.ID
		dm.Content = bm.Content
		dm.Role = bm.Role.Name()
		dm.Order = bm.Order
		dbMsgs[i] = dm
	}

	return dbMsgs
}
