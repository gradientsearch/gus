package conversationdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/gradientsearch/gus/business/domain/conversationbus"
	"github.com/gradientsearch/gus/business/sdk/sqldb"
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
	}
}

// NewWithTx constructs a new Store value replacing the sqlx DB
// value with a sqlx DB value that is currently inside a transaction.
func (s *Store) NewWithTx(tx sqldb.CommitRollbacker) (conversationbus.Storer, error) {
	ec, err := sqldb.GetExtContext(tx)
	if err != nil {
		return nil, err
	}

	store := Store{
		log: s.log,
		db:  ec,
	}

	return &store, nil
}

// QueryByID gets the specified conversation from the database.
func (s *Store) QueryById(ctx context.Context, userID uuid.UUID, conID uuid.UUID) (conversationbus.Conversation, error) {
	data := struct {
		ConversationID string `db:"conversation_id"`
		UserID         string `db:"user_id"`
	}{
		UserID:         userID.String(),
		ConversationID: conID.String(),
	}

	const q = `SELECT
    c.conversation_id,
    c.parent_message_id,
    c.user_id
FROM
    conversations c
WHERE
    c.user_id = :user_id
    AND c.conversation_id = :conversation_id
`

	var db conversation
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &db); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return conversationbus.Conversation{}, fmt.Errorf("db: %w", conversationbus.ErrNotFound)
		}
		return conversationbus.Conversation{}, fmt.Errorf("db: %w", err)
	}

	return toBusConversation(db)
}

func (s *Store) Create(ctx context.Context, c conversationbus.NewConversation) error {
	const q = `
	INSERT INTO conversations
		(conversation_id, user_id)
	VALUES
		(:conversation_id, :user_id)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDbConversation(c)); err != nil {
		return fmt.Errorf("db: %w", err)
	}

	return nil
}
