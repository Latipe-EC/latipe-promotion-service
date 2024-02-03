package storeserv

import (
	"context"
	"github.com/stretchr/testify/mock"
	"latipe-promotion-services/internal/adapter/storeserv/dto"
)

type StoreServiceMock struct {
	mock.Mock
}

func (s *StoreServiceMock) GetStoreByUserId(ctx context.Context, req *dto.GetStoreIdByUserRequest) (*dto.GetStoreIdByUserResponse, error) {
	args := s.Called(ctx, req)
	return args.Get(0).(*dto.GetStoreIdByUserResponse), args.Error(1)
}
