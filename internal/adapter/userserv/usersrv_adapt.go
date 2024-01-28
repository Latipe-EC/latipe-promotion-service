package userserv

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2/log"
	"latipe-promotion-services/config"
	"latipe-promotion-services/internal/adapter/userserv/dto"
	"latipe-promotion-services/pkgs/mapper"
)

type UserService struct {
	restyClient *resty.Client
	cfg         *config.Config
}

func NewUserService(cfg *config.Config) *UserService {
	restyClient := resty.New().SetDebug(true)

	return &UserService{
		restyClient: restyClient,
		cfg:         cfg,
	}
}

func (us UserService) GetAddressById(ctx context.Context, request *dto.GetAddressRequest) (*dto.GetAddressResponse, error) {
	resp, err := us.restyClient.
		SetBaseURL(us.cfg.AdapterService.UserService.UserURL).
		R().
		SetContext(ctx).
		SetDebug(false).
		Get(request.URL() + fmt.Sprintf("/%v", request.AddressId))

	if err != nil {
		log.Errorf("[%s] [Get address]: %s", "ERROR", err)
		return nil, err
	}

	if resp.StatusCode() >= 500 {
		log.Errorf("[%s] [Get address]: %s", "ERROR", resp.Body())
		return nil, errors.New("get address internal")
	}

	var regResp *dto.GetAddressResponse
	err = mapper.BindingStruct(resp.Body(), &regResp)
	if err != nil {
		log.Errorf("[%s] [Get product]: %s", "ERROR", err)
		return nil, err
	}

	return regResp, nil
}

func (us UserService) Authorization(ctx context.Context, req *dto.AuthorizationRequest) (*dto.AuthorizationResponse, error) {
	resp, err := us.restyClient.
		SetBaseURL(us.cfg.AdapterService.UserService.AuthURL).
		R().
		SetBody(req).
		SetContext(ctx).
		SetDebug(false).
		Post(req.URL())

	if err != nil {
		log.Errorf("[Authorize token]: %s", err)
		return nil, err
	}

	if resp.StatusCode() >= 500 {
		log.Errorf("[Authorize token]: %s", resp.Body())
		return nil, err
	}

	var regResp *dto.AuthorizationResponse

	if err := json.Unmarshal(resp.Body(), &regResp); err != nil {
		log.Errorf("[Authorize token]: %s", err)
		return nil, err
	}

	return regResp, nil
}
