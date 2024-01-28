package services

import (
	"github.com/google/wire"
	"latipe-promotion-services/internal/services/voucherserv"
)

var Set = wire.NewSet(voucherserv.NewVoucherService)
