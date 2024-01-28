package publisher

import (
	"github.com/google/wire"
	"latipe-promotion-services/internal/publisher/purchaseCreate"
)

var Set = wire.NewSet(purchaseCreate.NewReplyPurchaseTransactionPub)
