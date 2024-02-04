package vouchergrpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"latipe-promotion-services/internal/domain/dto"
	"latipe-promotion-services/internal/services/voucherserv"
	"latipe-promotion-services/pkgs/mapper"
	responses "latipe-promotion-services/pkgs/response"
	"latipe-promotion-services/pkgs/valid"
)

type voucherServer struct {
	voucherService *voucherserv.VoucherService
	UnimplementedVoucherServiceServer
}

func NewVoucherServerGRPC(voucherServ *voucherserv.VoucherService) VoucherServiceServer {
	return &voucherServer{
		voucherService: voucherServ,
	}
}

func (v voucherServer) CheckUsingVouchersForCheckout(ctx context.Context, request *CheckoutVoucherRequest) (*CheckoutVoucherResponse, error) {
	req := dto.PurchaseVoucherRequest{}
	var response CheckoutVoucherResponse

	if err := mapper.BindingStruct(request, &req); err != nil {
		return nil, err
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("%v", err))
	}

	resp, err := v.voucherService.CheckoutVoucher(ctx, &req)
	if err != nil {
		log.Errorf("%v", err)
		if errors.Is(err, responses.ErrUnableApplyVoucher) {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%v", err))
		}

		if errors.Is(err, responses.ErrVoucherExpiredOrOutOfStock) {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%v", err))
		}

		return nil, status.Error(codes.Internal, fmt.Sprintf("%v", err))
	}

	if err := mapper.BindingStructGrpc(resp, &response); err != nil {
		log.Errorf("%v", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("%v", err))
	}

	return &response, nil
}

func (v voucherServer) ApplyVoucherToPurchase(ctx context.Context, request *ApplyVoucherRequest) (*ApplyVoucherResponse, error) {
	req := dto.ApplyVoucherRequest{}
	response := ApplyVoucherResponse{}

	if err := mapper.BindingStruct(request, &req); err != nil {
		return nil, err
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("%v", err))
	}

	err := v.voucherService.UsingVoucherToOrder(ctx, &req)
	if err != nil {
		log.Errorf("%v", err)
		if errors.Is(err, responses.ErrUnableApplyVoucher) {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%v", err))
		}

		if errors.Is(err, responses.ErrVoucherExpiredOrOutOfStock) {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%v", err))
		}

		return nil, status.Error(codes.Internal, fmt.Sprintf("%v", err))
	}

	response.IsSuccess = true

	return &response, nil
}
