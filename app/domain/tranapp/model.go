package tranapp

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/gradientsearch/gus/app/sdk/errs"
)

// NewTran represents an example of cross domain transaction at the
// application layer.
type NewTran struct {
	Dialog NewDialog `json:"dialog"`
}

// Validate checks the data in the model is considered clean.
func (app NewTran) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

// Decode implements the decoder interface.
func (app *NewTran) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app NewDialog) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

// Decode implements the decoder interface.
func (app *NewDialog) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

type Message struct {
	ID      string `json:"id"`
	Role    string `json:"role"`
	Content string `json:"content"`
}

type NewDialog struct {
	Messages []Message
	UserID   uuid.UUID
}
