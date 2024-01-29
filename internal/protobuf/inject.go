package protobuf

import (
	"github.com/google/wire"
	"latipe-promotion-services/internal/protobuf/vouchergrpc"
)

var Set = wire.NewSet(vouchergrpc.NewVoucherServerGRPC)
