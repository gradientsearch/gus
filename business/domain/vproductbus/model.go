package vproductbus

import (
	"time"

	"github.com/gradientsearch/gus/business/types/money"
	"github.com/gradientsearch/gus/business/types/name"
	"github.com/gradientsearch/gus/business/types/quantity"
	"github.com/google/uuid"
)

// Product represents an individual product with extended information.
type Product struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        name.Name
	Cost        money.Money
	Quantity    quantity.Quantity
	DateCreated time.Time
	DateUpdated time.Time
	UserName    name.Name
}
