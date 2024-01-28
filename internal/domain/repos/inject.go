package repos

import "github.com/google/wire"

var Set = wire.NewSet(NewVoucherRepos)
