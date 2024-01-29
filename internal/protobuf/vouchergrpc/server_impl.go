package vouchergrpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"latipe-promotion-services/internal/domain/dto"
	"latipe-promotion-services/internal/services/voucherserv"
	"latipe-promotion-services/pkgs/mapper"
	"latipe-promotion-services/pkgs/valid"
)

type voucherServer struct {
	voucherService *voucherserv.VoucherService
	UnimplementedVoucherServiceGRPCServer
}

func NewVoucherServerGRPC(voucherServ *voucherserv.VoucherService) VoucherServiceGRPCServer {
	return &voucherServer{
		voucherService: voucherServ,
	}
}

func (v voucherServer) CheckingVoucher(ctx context.Context, request *CheckingVoucherRequest) (*CheckingVoucherResponse, error) {
	req := dto.CheckingVoucherRequest{}
	response := CheckingVoucherResponse{}

	if err := mapper.BindingStruct(request, &req); err != nil {
		return nil, err
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("%v", err))
	}

	resp, err := v.voucherService.CheckingVoucher(ctx, &req)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("%v", err))
	}

	if err := mapper.BindingStruct(resp, &response); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("%v", err))
	}

	return &response, nil
}

func (v voucherServer) ApplyVoucher(ctx context.Context, request *UseVoucherRequest) (*ApplyVoucherResponse, error) {
	req := dto.UseVoucherRequest{}
	response := ApplyVoucherResponse{}

	if err := mapper.BindingStruct(request, &req); err != nil {
		return nil, err
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("%v", err))
	}

	resp, err := v.voucherService.UseVoucher(ctx, &req)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("%v", err))
	}

	if err := mapper.BindingStruct(resp, &response); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("%v", err))
	}

	return &response, nil
}
