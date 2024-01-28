package adapter

import (
	"github.com/google/wire"
	"latipe-promotion-services/internal/adapter/userserv"
)

var Set = wire.NewSet(userserv.NewUserService)
