package chatdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/gradientsearch/gus/business/api/sqldb"
	"github.com/gradientsearch/gus/business/domain/chatbus"
	"github.com/gradientsearch/gus/business/domain/userbus"
	"github.com/gradientsearch/gus/foundation/logger"
	"github.com/jmoiron/sqlx"
)

// Store manages the set of APIs for user database access.
type Store struct {
	log *logger.Logger
	db  sqlx.ExtContext
}

// NewStore constructs the api for data access.
func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// QueryByID gets the specified conversation from the database.
func (s *Store) QueryById(ctx context.Context, userID uuid.UUID, conID uuid.UUID) (chatbus.Conversation, error) {
	data := struct {
		ConversationID string `db:"conversation_id"`
		UserID         string `db:"user_id"`
	}{
		UserID:         userID.String(),
		ConversationID: conID.String(),
	}

	const q = `	
SELECT
	c.id AS conversation_id,
	c.parent_message_id,
	m.id AS message_id,
	m.role,
	m.content
FROM
	conversations c
	JOIN messages m ON m.conversation_id = c.id;

WHERE
	c.user_id = :user_id
	AND c.id = :conversation_id;`

	var dbCon conversation
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbCon); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return chatbus.Conversation{}, fmt.Errorf("db: %w", userbus.ErrNotFound)
		}
		return chatbus.Conversation{}, fmt.Errorf("db: %w", err)
	}

	return toBusConversation(dbCon)
}
