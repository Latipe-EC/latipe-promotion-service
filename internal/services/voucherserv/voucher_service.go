package voucherserv

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	dto2 "latipe-promotion-services/internal/domain/dto"
	"latipe-promotion-services/internal/domain/entities"
	"latipe-promotion-services/internal/domain/message"
	"latipe-promotion-services/internal/domain/repos"
	"latipe-promotion-services/pkgs/mapper"
	"latipe-promotion-services/pkgs/pagable"
	"latipe-promotion-services/pkgs/response"
	"strings"
	"time"
)

type VoucherService struct {
	voucherRepos *repos.VoucherRepository
}

func NewVoucherService(provinceRepo *repos.VoucherRepository) *VoucherService {
	return &VoucherService{
		voucherRepos: provinceRepo,
	}
}

func (sh VoucherService) CreateVoucher(ctx context.Context, req *dto2.CreateVoucherRequest) (string, error) {

	vRequired := entities.VoucherRequire{
		MinRequire:    req.VoucherRequire.MinRequire,
		PaymentMethod: req.VoucherRequire.PaymentMethod,
	}

	voucherDAO := entities.Voucher{
		VoucherCode:     strings.ToUpper(req.VoucherCode),
		VoucherType:     req.VoucherType,
		VoucherCounts:   req.VoucherCounts,
		Detail:          req.Detail,
		DiscountPercent: req.DiscountPercent,
		DiscountValue:   req.DiscountValue,
	}
	voucherDAO.VoucherRequire = &vRequired
	voucherDAO.OwnerVoucher = "ADMIN"

	started := ParseStringToTime(req.StatedTime)
	voucherDAO.StatedTime = started

	ended := ParseStringToTime(req.EndedTime)
	if ended != started {
		voucherDAO.EndedTime = ended
	}

	voucherDAO.Status = entities.PENDING

	resp, err := sh.voucherRepos.CreateVoucher(ctx, &voucherDAO)
	if err != nil {
		log.Errorf("%v", err)
		return "", err
	}

	return resp, err
}

func (sh VoucherService) GetVoucherByCode(ctx context.Context, code string) (*dto2.VoucherRespDetail, error) {
	voucher, err := sh.voucherRepos.GetByCode(ctx, strings.ToUpper(code))
	if err != nil {
		return nil, err
	}

	voucherResp := dto2.VoucherRespDetail{}
	if err := mapper.BindingStruct(voucher, &voucherResp); err != nil {
		return nil, err
	}

	return &voucherResp, err
}

func (sh VoucherService) UpdateVoucherStatus(ctx context.Context, req *dto2.UpdateVoucherRequest) error {
	voucher, err := sh.voucherRepos.GetByCode(ctx, strings.ToUpper(req.VoucherCode))
	if err != nil {
		return err
	}

	voucher.Status = req.Status

	if err := sh.voucherRepos.UpdateStatus(ctx, voucher); err != nil {
		return err
	}

	return nil
}

func (sh VoucherService) GetVoucherById(ctx context.Context, id string) (*dto2.VoucherRespDetail, error) {
	voucher, err := sh.voucherRepos.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	voucherResp := dto2.VoucherRespDetail{}
	if err := mapper.BindingStruct(voucher, &voucherResp); err != nil {
		return nil, err
	}

	return &voucherResp, err
}

func (sh VoucherService) GetAllVoucher(ctx context.Context, voucherCode string, query *pagable.Query) (*pagable.ListResponse, error) {
	vouchers, total, err := sh.voucherRepos.GetAll(ctx, voucherCode, query)
	if err != nil {
		return nil, err
	}

	var voucherResp []dto2.VoucherRespDetail
	if err := mapper.BindingStruct(vouchers, &voucherResp); err != nil {
		return nil, err
	}

	listResp := pagable.ListResponse{}
	listResp.Items = voucherResp
	listResp.Page = query.GetPage()
	listResp.Size = query.GetSize()
	listResp.Total = query.GetTotalPages(int(total))
	listResp.HasMore = query.GetHasMore(int(total))

	return &listResp, err
}

func (sh VoucherService) GetUserVoucher(ctx context.Context, voucherCode string, query *pagable.Query) (*pagable.ListResponse, error) {
	vouchers, total, err := sh.voucherRepos.GetVoucherForUser(ctx, voucherCode, query)
	if err != nil {
		return nil, err
	}

	var voucherResp []dto2.VoucherUserDetail
	if err := mapper.BindingStruct(vouchers, &voucherResp); err != nil {
		return nil, err
	}

	listResp := pagable.ListResponse{}
	listResp.Items = voucherResp
	listResp.Page = query.GetPage()
	listResp.Size = query.GetSize()
	listResp.Total = query.GetTotalPages(int(total))
	listResp.HasMore = query.GetHasMore(int(total))

	return &listResp, err
}

func (sh VoucherService) UseVoucher(ctx context.Context, req *dto2.UseVoucherRequest) (*dto2.UseVoucherResponse, error) {
	var vouchers []*entities.Voucher

	for _, i := range req.Vouchers {
		voucher, err := sh.voucherRepos.GetByCode(ctx, i)
		if err != nil {
			return nil, err
		}

		if voucher.VoucherCounts > 0 {
			voucher.VoucherCounts -= 1
		} else {
			return nil, responses.ErrVoucherExpiredOrOutOfStock
		}

		if !voucher.EndedTime.After(time.Now()) && !voucher.StatedTime.Before(time.Now()) ||
			voucher.Status != entities.ACTIVE {
			return nil, responses.ErrVoucherExpiredOrOutOfStock
		}

		vouchers = append(vouchers, voucher)
	}

	if len(vouchers) != len(req.Vouchers) {
		return nil, responses.ErrBadRequest
	}

	if len(req.Vouchers) > 1 && vouchers[0].VoucherType == vouchers[1].VoucherType {
		return nil, responses.ErrDuplicateType
	}

	if err := sh.voucherRepos.UpdateVoucherCounts(ctx, vouchers); err != nil {
		return nil, err
	}

	for _, i := range vouchers {
		usingLog := entities.VoucherUsingLog{
			VoucherCode: i.VoucherCode,
			VoucherID:   i.ID,
			OrderID:     req.OrderUUID,
			Status:      1,
			CreatedAt:   time.Now(),
		}

		err := sh.voucherRepos.CreateLogUseVoucher(ctx, &usingLog)
		if err != nil {
			log.Errorf("%v", err)
		}
	}

	resp := dto2.UseVoucherResponse{}
	resp.IsSuccess = true
	if err := mapper.BindingStruct(vouchers, &resp.Items); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (sh VoucherService) CheckingVoucher(ctx context.Context, req *dto2.CheckingVoucherRequest) (*dto2.UseVoucherResponse, error) {
	var vouchers []*entities.Voucher

	for _, i := range req.Vouchers {
		voucher, err := sh.voucherRepos.GetByCode(ctx, i)
		if err != nil {
			return nil, err
		}

		if voucher.VoucherCounts > 0 {
			voucher.VoucherCounts -= 1
		} else {
			return nil, responses.ErrVoucherExpiredOrOutOfStock
		}

		if !voucher.EndedTime.After(time.Now()) && !voucher.StatedTime.Before(time.Now()) ||
			voucher.Status != entities.ACTIVE {
			return nil, responses.ErrVoucherExpiredOrOutOfStock
		}

		//check required
		if voucher.VoucherRequire != nil {
			if int64(req.OrderTotalAmount) < voucher.VoucherRequire.MinRequire {
				return nil, responses.ErrUnableApplyVoucher
			}

			if voucher.VoucherRequire.PaymentMethod != 0 && req.PaymentMethod != voucher.VoucherRequire.PaymentMethod {
				return nil, responses.ErrUnableApplyVoucher
			}
		}

		vouchers = append(vouchers, voucher)
	}

	if len(vouchers) != len(req.Vouchers) {
		return nil, responses.ErrNotFoundRecord
	}

	if len(req.Vouchers) > 1 && vouchers[0].VoucherType == vouchers[1].VoucherType {
		return nil, responses.ErrDuplicateType
	}

	resp := dto2.UseVoucherResponse{}
	resp.IsSuccess = true
	if err := mapper.BindingStruct(vouchers, &resp.Items); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (sh VoucherService) RollBackVoucher(ctx context.Context, req *dto2.UseVoucherRequest) error {
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

func (sh VoucherService) UseVoucherTransaction(ctx context.Context, req *message.CreatePurchaseMessage) error {
	var vouchers []*entities.Voucher

	for _, i := range req.Vouchers {
		voucher, err := sh.voucherRepos.GetByCode(ctx, i)
		if err != nil {
			return err
		}

		if voucher.VoucherCounts > 0 {
			voucher.VoucherCounts -= 1
		} else {
			return responses.ErrVoucherExpiredOrOutOfStock
		}

		if !voucher.EndedTime.After(time.Now()) && !voucher.StatedTime.Before(time.Now()) ||
			voucher.Status != entities.ACTIVE {
			return responses.ErrVoucherExpiredOrOutOfStock
		}

		vouchers = append(vouchers, voucher)
	}

	if len(vouchers) != len(req.Vouchers) {
		return responses.ErrBadRequest
	}

	if len(req.Vouchers) > 1 && vouchers[0].VoucherType == vouchers[1].VoucherType {
		return responses.ErrDuplicateType
	}

	if err := sh.voucherRepos.UpdateVoucherCounts(ctx, vouchers); err != nil {
		return err
	}

	for _, i := range vouchers {
		usingLog := entities.VoucherUsingLog{
			VoucherCode: i.VoucherCode,
			VoucherID:   i.ID,
			OrderID:     req.OrderID,
			Status:      1,
			CreatedAt:   time.Now(),
		}

		err := sh.voucherRepos.CreateLogUseVoucher(ctx, &usingLog)
		if err != nil {
			log.Errorf("%v", err)
		}
	}

	resp := dto2.UseVoucherResponse{}
	resp.IsSuccess = true
	if err := mapper.BindingStruct(vouchers, &resp.Items); err != nil {
		return err
	}

	return nil
}

func (sh VoucherService) RollbackVoucherTransaction(ctx context.Context, req *message.RollbackPurchaseMessage) error {
	return nil
}
