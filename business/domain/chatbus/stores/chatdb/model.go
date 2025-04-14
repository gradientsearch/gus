package chatdb

import "github.com/gradientsearch/gus/business/domain/chatbus"

type conversation struct {
}

func toBusConversation(dbCon conversation) (chatbus.Conversation, error) {

	return chatbus.Conversation{}, nil
}
