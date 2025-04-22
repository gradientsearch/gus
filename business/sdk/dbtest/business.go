package dbtest

import (
	"time"

	"github.com/gradientsearch/gus/business/domain/userbus"
	"github.com/gradientsearch/gus/business/domain/userbus/stores/usercache"
	"github.com/gradientsearch/gus/business/domain/userbus/stores/userdb"

	"github.com/gradientsearch/gus/business/sdk/delegate"
	"github.com/gradientsearch/gus/foundation/logger"
	"github.com/jmoiron/sqlx"
)

// BusDomain represents all the business domain apis needed for testing.
type BusDomain struct {
	Delegate *delegate.Delegate
	User     *userbus.Business
}

func newBusDomains(log *logger.Logger, db *sqlx.DB) BusDomain {
	delegate := delegate.New(log)
	userBus := userbus.NewBusiness(log, delegate, usercache.NewStore(log, userdb.NewStore(log, db), time.Hour))

	return BusDomain{
		Delegate: delegate,
		User:     userBus,
	}
}
