package chatdb

import (
	"github.com/google/uuid"
	"github.com/gradientsearch/gus/business/domain/chatbus"
)

type conversation struct {
	ID              uuid.UUID `db:"conversation_id"`
	ParentMessageID uuid.UUID `db:"parent_message_id"`
	UserID          uuid.UUID `db:"user_id"`
}

type conversationMessage struct {
	conversation
	message
}

type message struct {
	ID      uuid.UUID `db:"message_id"`
	Role    string    `db:"role"`
	Content string    `db:"content"`
	Order   int       `db:"order"`
}

func toBusConversation(dbCon []conversationMessage) (chatbus.Conversation, error) {

	return chatbus.Conversation{}, nil
}

func toDbConversation(busCon chatbus.Conversation) conversation {
	var dbCon = conversation{}
	dbCon.ID = busCon.ID
	dbCon.ParentMessageID = busCon.ParentMessageID
	dbCon.UserID = busCon.UserID
	return dbCon
}

func toDbMessages(busMsgs []chatbus.Message) []message {
	dbMsgs := make([]message, len(busMsgs))
	for i, bm := range busMsgs {
		dm := message{}
		dm.ID = bm.ID
		dm.Content = bm.Content
		dm.Role = bm.Role.Name()
		dm.Order = bm.Order
		dbMsgs[i] = dm
	}

	return dbMsgs
}
