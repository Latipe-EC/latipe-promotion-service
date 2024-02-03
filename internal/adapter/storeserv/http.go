package storeserv

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/wire"
	"latipe-promotion-services/config"
	"latipe-promotion-services/internal/adapter/storeserv/dto"
	responses "latipe-promotion-services/pkgs/response"
)

var Set = wire.NewSet(
	NewStoreServiceAdapter,
)

type httpAdapter struct {
	client *resty.Client
}

func NewStoreServiceAdapter(config *config.Config) StoreService {
	restyClient := resty.New()
	restyClient.
		SetBaseURL(config.AdapterService.StoreService.BaseURL).
		SetHeader("X-INTERNAL-SERVICE", config.AdapterService.StoreService.InternalKey)
	return httpAdapter{
		client: restyClient,
	}
}

func (h httpAdapter) GetStoreByUserId(ctx context.Context, req *dto.GetStoreIdByUserRequest) (*dto.GetStoreIdByUserResponse, error) {
	resp, err := h.client.R().
		SetContext(ctx).
		SetHeader("Authorization", fmt.Sprintf("Bearer %v", req.BaseHeader.BearToken)).
		Get(req.URL() + req.UserID)

	if err != nil {
		log.Errorf("[Get store]: %s", err)
		return nil, err
	}

	if resp.StatusCode() >= 500 {
		log.Errorf("[Get store]: %s", resp.Body())
		return nil, responses.ErrInternalServer
	}

	if resp.StatusCode() >= 400 {
		log.Errorf("[Get store]: %s", resp.Body())
		return nil, responses.ErrBadRequest
	}

	regResp := dto.GetStoreIdByUserResponse{
		StoreID: resp.String(),
	}

	return &regResp, nil
}
