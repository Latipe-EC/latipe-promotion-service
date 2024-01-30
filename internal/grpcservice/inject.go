package grpcservice

import (
	"github.com/google/wire"
	"latipe-promotion-services/internal/grpcservice/interceptor"
	"latipe-promotion-services/internal/grpcservice/vouchergrpc"
)

var Set = wire.NewSet(
	vouchergrpc.NewVoucherServerGRPC,
	interceptor.NewGrpcInterceptor,
)
