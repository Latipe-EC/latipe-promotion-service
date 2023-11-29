package voucherserv

import (
	"context"
	"errors"
	"latipe-promotion-services/domain/dto"
	entities "latipe-promotion-services/domain/entities"
	repos "latipe-promotion-services/domain/repos"
	"latipe-promotion-services/pkgs/mapper"
	"latipe-promotion-services/pkgs/pagable"
	"time"
)

type VoucherService struct {
	voucherRepos *repos.VoucherRepository
}

func NewVoucherService(provinceRepo *repos.VoucherRepository) VoucherService {
	return VoucherService{
		voucherRepos: provinceRepo,
	}
}

func (sh VoucherService) CreateVoucher(ctx context.Context, req *dto.CreateVoucherRequest) (string, error) {

	voucherDAO := entities.Voucher{
		VoucherCode:     req.VoucherCode,
		VoucherType:     req.VoucherType,
		VoucherCounts:   req.VoucherCounts,
		Detail:          req.Detail,
		DiscountPercent: req.DiscountPercent,
		DiscountValue:   req.DiscountValue,
		VoucherRequire: entities.VoucherRequire{
			MinRequire:    req.VoucherRequire.MinRequire,
			MemberType:    req.VoucherRequire.MemberType,
			PaymentMethod: req.VoucherRequire.PaymentMethod,
		},
	}

	req.OwnerVoucherId = voucherDAO.OwnerVoucher
	if req.OwnerVoucherId != "ADMIN" {
		voucherDAO.VoucherRequire.RequiredOwnerProdId = req.OwnerVoucherId
	}
	started := ParseStringToTime(req.StatedTime)
	voucherDAO.StatedTime = started

	ended := ParseStringToTime(req.EndedTime)
	if ended != started {
		voucherDAO.EndedTime = ended
	}

	voucherDAO.Status = entities.PENDING

	resp, err := sh.voucherRepos.CreateVoucher(ctx, &voucherDAO)
	if err != nil {
		return "", err
	}

	return resp, err
}

func (sh VoucherService) GetVoucherByCode(ctx context.Context, code string) (*dto.VoucherRespDetail, error) {
	voucher, err := sh.voucherRepos.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	voucherResp := dto.VoucherRespDetail{}
	if err := mapper.BindingStruct(voucher, &voucherResp); err != nil {
		return nil, err
	}

	return &voucherResp, err
}

func (sh VoucherService) GetVoucherById(ctx context.Context, id string) (*dto.VoucherRespDetail, error) {
	voucher, err := sh.voucherRepos.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	voucherResp := dto.VoucherRespDetail{}
	if err := mapper.BindingStruct(voucher, &voucherResp); err != nil {
		return nil, err
	}

	return &voucherResp, err
}

func (sh VoucherService) GetAllVoucher(ctx context.Context, query *pagable.Query) (*pagable.ListResponse, error) {
	vouchers, err := sh.voucherRepos.GetAll(ctx, query)
	if err != nil {
		return nil, err
	}

	total, _ := sh.voucherRepos.Total(ctx, query)

	var voucherResp []dto.VoucherRespDetail
	if err := mapper.BindingStruct(vouchers, &voucherResp); err != nil {
		return nil, err
	}

	listResp := pagable.ListResponse{}
	listResp.Data = voucherResp
	listResp.Page = query.GetPage()
	listResp.Size = query.GetSize()
	listResp.Total = query.GetTotalPages(int(total))
	listResp.HasMore = query.GetHasMore(int(total))

	return &listResp, err
}

func (sh VoucherService) GetUserVoucher(ctx context.Context, query *pagable.Query) (*pagable.ListResponse, error) {
	vouchers, err := sh.voucherRepos.GetVoucherForUser(ctx, query)
	if err != nil {
		return nil, err
	}

	total, _ := sh.voucherRepos.Total(ctx, query)

	var voucherResp []dto.VoucherUserDetail
	if err := mapper.BindingStruct(vouchers, &voucherResp); err != nil {
		return nil, err
	}

	listResp := pagable.ListResponse{}
	listResp.Data = voucherResp
	listResp.Page = query.GetPage()
	listResp.Size = query.GetSize()
	listResp.Total = query.GetTotalPages(int(total))
	listResp.HasMore = query.GetHasMore(int(total))

	return &listResp, err
}

func (sh VoucherService) UseVoucher(ctx context.Context, req *dto.UseVoucherRequest) (*dto.UseVoucherResponse, error) {
	var vouchers []*entities.Voucher

	for _, i := range req.Vouchers {
		voucher, err := sh.voucherRepos.GetByCode(ctx, i)
		if err != nil {
			return nil, err
		}

		if voucher.VoucherCounts > 0 {
			voucher.VoucherCounts -= 1
		} else {
			return nil, errors.New("the voucher amount was sold out")
		}

		if !voucher.EndedTime.After(time.Now()) && !voucher.StatedTime.Before(time.Now()) &&
			voucher.Status == entities.INACTIVE {
			return nil, errors.New("the voucher has expired")
		}

		vouchers = append(vouchers, voucher)
	}

	if len(vouchers) != len(req.Vouchers) {
		return nil, errors.New("voucher code was not found")
	}

	if len(req.Vouchers) > 1 && vouchers[0].VoucherType == vouchers[1].VoucherType {
		return nil, errors.New("just one for one type")
	}

	if err := sh.voucherRepos.UpdateVoucherCounts(ctx, vouchers); err != nil {
		return nil, err
	}

	resp := dto.UseVoucherResponse{}
	resp.IsSuccess = true
	if err := mapper.BindingStruct(vouchers, &resp.Items); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (sh VoucherService) RollBackVoucher(ctx context.Context, req *dto.UseVoucherRequest) error {
	var vouchers []*entities.Voucher

	for _, i := range req.Vouchers {
		voucher, err := sh.voucherRepos.GetByCode(ctx, i)
		if err != nil {
			return err
		}

		voucher.VoucherCounts += 1
		vouchers = append(vouchers, voucher)
	}

	if err := sh.voucherRepos.UpdateVoucherCounts(ctx, vouchers); err != nil {
		return err
	}

	return nil
}
