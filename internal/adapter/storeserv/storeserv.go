package storeserv

import (
	"context"
	"latipe-promotion-services/internal/adapter/storeserv/dto"
)

type StoreService interface {
	GetStoreByUserId(ctx context.Context, req *dto.GetStoreIdByUserRequest) (*dto.GetStoreIdByUserResponse, error)
}
