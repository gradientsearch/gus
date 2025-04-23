package dialogdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/gradientsearch/gus/business/domain/dialogbus"
	"github.com/gradientsearch/gus/business/sdk/sqldb"
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

// NewWithTx constructs a new Store value replacing the sqlx DB
// value with a sqlx DB value that is currently inside a transaction.
func (s *Store) NewWithTx(tx sqldb.CommitRollbacker) (dialogbus.Storer, error) {
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
func (s *Store) QueryById(ctx context.Context, userID uuid.UUID, conID uuid.UUID) (dialogbus.Dialog, error) {
	data := struct {
		ConversationID string `db:"conversation_id"`
		UserID         string `db:"user_id"`
	}{
		UserID:         userID.String(),
		ConversationID: conID.String(),
	}

	const q = `SELECT
    m.conversation_id,
    m.user_id,
    m.message_id,
    m.role,
    m.content,
    m.order
FROM
  	messages m 
WHERE
    m.user_id = :user_id
    AND m.conversation_id = :conversation_id
ORDER BY
    m.order ASC;
`

	var dbMessages []message
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &dbMessages); err != nil {
		if !errors.Is(err, sqldb.ErrDBNotFound) {
			return dialogbus.Dialog{}, fmt.Errorf("db: %w", dialogbus.ErrQuery)
		}
		dbMessages = []message{}
	}

	return toBusDialog(dbMessages)
}

// Create
func (s *Store) Create(ctx context.Context, bus dialogbus.Dialog) error {
	const q = `
	INSERT INTO messages
		(message_id, conversation_id, user_id, role, content, "order")
	VALUES
		(:message_id, :conversation_id, :user_id, :role, :content, :order)`

	for _, m := range toDbMessages(bus) {
		if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, m); err != nil {
			return fmt.Errorf("db: %w", err)
		}
	}

	return nil
}
