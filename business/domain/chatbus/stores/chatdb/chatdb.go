package chatdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/gradientsearch/gus/business/api/sqldb"
	"github.com/gradientsearch/gus/business/domain/chatbus"
	"github.com/gradientsearch/gus/foundation/logger"
	"github.com/jmoiron/sqlx"
)

// Store manages the set of APIs for user database access.
type Store struct {
	log *logger.Logger
	db  sqlx.ExtContext
	tx  *sqlx.DB
}

// NewStore constructs the api for data access.
func NewStore(log *logger.Logger, db *sqlx.DB) *Store {

	return &Store{
		log: log,
		db:  db,
		tx:  db,
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
			return chatbus.Conversation{}, fmt.Errorf("db: %w", chatbus.ErrNotFound)
		}
		return chatbus.Conversation{}, fmt.Errorf("db: %w", err)
	}

	return toBusConversation(dbCon)
}

func (s *Store) Create(ctx context.Context, c chatbus.Conversation) error {
	tx, err := s.tx.DB.Begin()
	if err != nil {
		return fmt.Errorf("db: %w", err)
	}
	const q1 = `
	INSERT INTO conversations
		(conversation_id, parent_message_id, user_id)
	VALUES
		($1, $2, $3)`

	if _, err = tx.ExecContext(ctx, q1, c.ID, c.ParentMessageID, c.UserID); err != nil {
		return fmt.Errorf("db: %w", err)
	}

	const q2 = `
	INSERT INTO messages
		(message_id, conversation_id, user_id, role, content, order)
	VALUES
		($1, $2, $3, $4, $5, $6)`

	for _, m := range toDbMessages(c.Messages) {
		if _, err = tx.ExecContext(ctx, q2, m.ID, c.ID, c.UserID, m.Role, m.Content, m.Order); err != nil {
			return fmt.Errorf("db: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("db: %w", err)
	}

	return nil
}
func (s *Store) Update(ctx context.Context, c chatbus.Conversation) error {
	return nil
}
