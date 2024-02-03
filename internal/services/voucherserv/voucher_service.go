package voucherserv

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	dto "latipe-promotion-services/internal/domain/dto"
	"latipe-promotion-services/internal/domain/entities"
	"latipe-promotion-services/internal/domain/message"
	"latipe-promotion-services/internal/domain/repos"
	"latipe-promotion-services/internal/publisher/purchaseCreate"
	"latipe-promotion-services/pkgs/mapper"
	"latipe-promotion-services/pkgs/pagable"
	responses "latipe-promotion-services/pkgs/response"
	"strings"
	"time"
)

type VoucherService struct {
	voucherRepos *repos.VoucherRepository
	voucherReply *purchaseCreate.ReplyPurchaseTransactionPub
}

func NewVoucherService(provinceRepo *repos.VoucherRepository,
	voucherReply *purchaseCreate.ReplyPurchaseTransactionPub) *VoucherService {
	return &VoucherService{
		voucherRepos: provinceRepo,
		voucherReply: voucherReply,
	}
}

func (sh VoucherService) CreateVoucher(ctx context.Context, req *dto.CreateVoucherRequest) (string, error) {

	var dao entities.Voucher

	if err := ValidateVoucherRequest(req); err != nil {
		return "", err
	}

	switch req.VoucherType {
	case entities.FREE_SHIP:
		vRequired := entities.VoucherRequire{
			MinRequire:        req.VoucherRequire.MinRequire,
			PaymentMethod:     req.VoucherRequire.PaymentMethod,
			MaxVoucherPerUser: req.VoucherRequire.MaxVoucherPerUser,
		}

		discount := entities.DiscountData{
			ShippingValue: req.DiscountData.ShippingValue,
		}

		voucherDAO := entities.Voucher{
			VoucherCode:    strings.ToUpper(req.VoucherCode),
			VoucherType:    req.VoucherType,
			VoucherCounts:  req.VoucherCounts,
			Detail:         req.Detail,
			StatedTime:     ParseStringToTime(req.StatedTime),
			EndedTime:      ParseStringToTime(req.EndedTime),
			OwnerVoucher:   "ADMIN",
			Status:         entities.PENDING,
			VoucherRequire: &vRequired,
			DiscountData:   &discount,
		}
		dao = voucherDAO
	case entities.PAYMENT_DISCOUNT:
		vRequired := entities.VoucherRequire{
			MinRequire:        req.VoucherRequire.MinRequire,
			PaymentMethod:     req.VoucherRequire.PaymentMethod,
			MaxVoucherPerUser: req.VoucherRequire.MaxVoucherPerUser,
		}

		discount := entities.DiscountData{
			DiscountType:    req.DiscountData.DiscountType,
			DiscountValue:   req.DiscountData.DiscountValue,
			DiscountPercent: req.DiscountData.DiscountPercent,
			MaximumValue:    req.DiscountData.MaximumValue,
		}

		voucherDAO := entities.Voucher{
			VoucherCode:    strings.ToUpper(req.VoucherCode),
			VoucherType:    req.VoucherType,
			VoucherCounts:  req.VoucherCounts,
			Detail:         req.Detail,
			StatedTime:     ParseStringToTime(req.StatedTime),
			EndedTime:      ParseStringToTime(req.EndedTime),
			OwnerVoucher:   "ADMIN",
			Status:         entities.PENDING,
			VoucherRequire: &vRequired,
			DiscountData:   &discount,
		}
		dao = voucherDAO
	case entities.STORE_DISCOUNT:
		vRequired := entities.VoucherRequire{
			MinRequire:        req.VoucherRequire.MinRequire,
			PaymentMethod:     req.VoucherRequire.PaymentMethod,
			MaxVoucherPerUser: req.VoucherRequire.MaxVoucherPerUser,
		}

		discount := entities.DiscountData{
			DiscountType:    req.DiscountData.DiscountType,
			DiscountValue:   req.DiscountData.DiscountValue,
			DiscountPercent: req.DiscountData.DiscountPercent,
			MaximumValue:    req.DiscountData.MaximumValue,
		}

		voucherDAO := entities.Voucher{
			VoucherCode:    strings.ToUpper(req.VoucherCode),
			VoucherType:    req.VoucherType,
			VoucherCounts:  req.VoucherCounts,
			Detail:         req.Detail,
			StatedTime:     ParseStringToTime(req.StatedTime),
			EndedTime:      ParseStringToTime(req.EndedTime),
			OwnerVoucher:   req.StoreID,
			Status:         entities.PENDING,
			VoucherRequire: &vRequired,
			DiscountData:   &discount,
		}
		dao = voucherDAO
	}

	resp, err := sh.voucherRepos.CreateVoucher(ctx, &dao)
	if err != nil {
		log.Errorf("%v", err)
		return "", err
	}

	return resp, err
}

func (sh VoucherService) GetVoucherByCode(ctx context.Context, code string) (*dto.VoucherRespDetail, error) {
	voucher, err := sh.voucherRepos.GetByCode(ctx, strings.ToUpper(code))
	if err != nil {
		return nil, err
	}

	voucherResp := dto.VoucherRespDetail{}
	if err := mapper.BindingStruct(voucher, &voucherResp); err != nil {
		return nil, err
	}

	return &voucherResp, err
}

func (sh VoucherService) UpdateVoucherStatus(ctx context.Context, req *dto.UpdateVoucherRequest) error {
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

func (sh VoucherService) GetAllVoucher(ctx context.Context, voucherCode string, query *pagable.Query) (*pagable.ListResponse, error) {
	vouchers, total, err := sh.voucherRepos.GetAll(ctx, voucherCode, query)
	if err != nil {
		return nil, err
	}

	var voucherResp []dto.VoucherRespDetail
	if err := mapper.BindingStruct(vouchers, &voucherResp); err != nil {
		return nil, err
	}

	listResp := pagable.ListResponse{}
	listResp.Items = voucherResp
	listResp.Page = query.GetPage()
	listResp.Size = query.GetSize()
	listResp.Total = query.GetTotalPages(total)
	listResp.HasMore = query.GetHasMore(total)

	return &listResp, err
}

func (sh VoucherService) GetUserVoucher(ctx context.Context, voucherCode string, userId string, query *pagable.Query) (*pagable.ListResponse, error) {
	vouchers, total, err := sh.voucherRepos.GetVoucherForUser(ctx, voucherCode, query)
	if err != nil {
		return nil, err
	}

	var voucherResp []dto.VoucherUserDetail
	if err := mapper.BindingStruct(vouchers, &voucherResp); err != nil {
		return nil, err
	}

	for index, i := range voucherResp {
		totalVouchers, err := sh.voucherRepos.CheckUsableVoucherOfUser(ctx, voucherCode, userId)
		if err != nil {
			return nil, err
		}

		if totalVouchers >= i.VoucherRequire.MaxVoucherPerUser {
			voucherResp[index].Usable = false
		} else {
			voucherResp[index].Usable = true
		}
	}

	listResp := pagable.ListResponse{}
	listResp.Items = voucherResp
	listResp.Page = query.GetPage()
	listResp.Size = query.GetSize()
	listResp.Total = query.GetTotalPages(total)
	listResp.HasMore = query.GetHasMore(total)

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
			OrderID:     req.OrderID,
			Status:      1,
			CreatedAt:   time.Now(),
		}

		err := sh.voucherRepos.CreateLogUseVoucher(ctx, &usingLog)
		if err != nil {
			log.Errorf("%v", err)
		}
	}

	resp := dto.UseVoucherResponse{}
	resp.IsSuccess = true
	if err := mapper.BindingStruct(vouchers, &resp.Items); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (sh VoucherService) CheckingVoucher(ctx context.Context, req *dto.CheckingVoucherRequest) (*dto.UseVoucherResponse, error) {
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

func (sh VoucherService) UseVoucherTransaction(ctx context.Context, msg *message.CreatePurchaseMessage) error {
	// Initialize reply message
	msgReply := message.ReplyPurchaseMessage{
		OrderID: msg.OrderID,
		Status:  message.COMMIT_SUCCESS, // Assume success until an error occurs
	}

	// Defer the reply message publishing
	defer func(replyPub *purchaseCreate.ReplyPurchaseTransactionPub, replyMsg *message.ReplyPurchaseMessage) {
		err := replyPub.ReplyPurchaseMessage(replyMsg)
		if err != nil {
			log.Error(err)
		}
	}(sh.voucherReply, &msgReply)

	var vouchers []*entities.Voucher

	for _, voucherCode := range msg.Vouchers {
		voucher, err := sh.voucherRepos.GetByCode(ctx, voucherCode)
		if err != nil {
			return err
		}

		// Check voucher validity
		if voucher.VoucherCounts <= 0 || voucher.Status != entities.ACTIVE ||
			!voucher.StatedTime.Before(time.Now()) || !voucher.EndedTime.After(time.Now()) {
			return responses.ErrVoucherExpiredOrOutOfStock
		}

		// Decrement voucher counts
		voucher.VoucherCounts--

		vouchers = append(vouchers, voucher)
	}

	// Check if all vouchers were found
	if len(vouchers) != len(msg.Vouchers) {
		msgReply.Status = message.COMMIT_FAIL
		return responses.ErrBadRequest
	}

	// Check for duplicate voucher types
	if len(msg.Vouchers) > 1 && vouchers[0].VoucherType == vouchers[1].VoucherType {
		msgReply.Status = message.COMMIT_FAIL
		return responses.ErrDuplicateType
	}

	// Update voucher counts in the repository
	if err := sh.voucherRepos.UpdateVoucherCounts(ctx, vouchers); err != nil {
		msgReply.Status = message.COMMIT_FAIL
		return err
	}

	// Create voucher using logs
	for _, voucher := range vouchers {
		usingLog := entities.VoucherUsingLog{
			VoucherCode: voucher.VoucherCode,
			VoucherID:   voucher.ID,
			OrderID:     msg.OrderID,
			Status:      1,
			CreatedAt:   time.Now(),
		}

		if err := sh.voucherRepos.CreateLogUseVoucher(ctx, &usingLog); err != nil {
			// Log the error, but continue processing other vouchers
			log.Errorf("Error creating voucher usage log: %v", err)
			msgReply.Status = message.COMMIT_FAIL
		}
	}

	return nil
}

func (sh VoucherService) RollbackVoucherTransaction(ctx context.Context, req *message.RollbackPurchaseMessage) error {
	return nil
}

// Store
func (sh VoucherService) StoreCancelVoucher(ctx context.Context, req *dto.CancelUpdateVoucherRequest) error {
	voucher, err := sh.voucherRepos.GetByCode(ctx, strings.ToUpper(req.VoucherCode))
	if err != nil {
		return err
	}

	if voucher.Status != entities.ACTIVE {
		return responses.ErrPermissionDenied
	}
	voucher.Status = entities.INACTIVE

	if err := sh.voucherRepos.UpdateStatus(ctx, voucher); err != nil {
		return err
	}

	return nil
}

func (sh VoucherService) GetAllVoucherOfStore(ctx context.Context, voucherCode string, storeId string, query *pagable.Query) (*pagable.ListResponse, error) {
	vouchers, total, err := sh.voucherRepos.GetVoucherOfStore(ctx, storeId, voucherCode, query)
	if err != nil {
		return nil, err
	}

	var voucherResp []dto.VoucherRespDetail
	if err := mapper.BindingStruct(vouchers, &voucherResp); err != nil {
		return nil, err
	}

	listResp := pagable.ListResponse{}
	listResp.Items = voucherResp
	listResp.Page = query.GetPage()
	listResp.Size = query.GetSize()
	listResp.Total = query.GetTotalPages(total)
	listResp.HasMore = query.GetHasMore(total)

	return &listResp, err
}
