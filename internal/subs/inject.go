package subscriber

import (
	"github.com/google/wire"
	"latipe-promotion-services/internal/subs/createPurchase"
)

var Set = wire.NewSet(
	createPurchase.NewPurchaseCreateSubscriber,
	createPurchase.NewPurchaseRollbackSubscriber,
)
