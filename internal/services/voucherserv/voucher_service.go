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

func (sh VoucherService) CheckoutVoucher(ctx context.Context, req *dto.PurchaseVoucherRequest) (*dto.UseVoucherResponse, error) {
	resp := dto.UseVoucherResponse{}

	if req.PaymentVoucher != nil {
		voucher, err := sh.voucherRepos.GetByCode(ctx, req.PaymentVoucher.VoucherCode)
		if err != nil {
			return nil, err
		}

		if err := sh.validateVoucherRequirement(ctx, req, voucher); err != nil {
			return nil, err
		}

		if err := mapper.BindingStruct(voucher, &resp.PaymentVoucher); err != nil {
			return nil, err
		}
	}

	if req.FreeShippingVoucher != nil {
		voucher, err := sh.voucherRepos.GetByCode(ctx, req.FreeShippingVoucher.VoucherCode)
		if err != nil {
			return nil, err
		}

		if err := sh.validateVoucherRequirement(ctx, req, voucher); err != nil {
			return nil, err
		}

		if err := mapper.BindingStruct(voucher, &resp.ShippingVoucher); err != nil {
			return nil, err
		}
	}

	if len(req.ShopVouchers) > 0 {
		var storeVoucher []*entities.Voucher

		for _, i := range req.ShopVouchers {
			voucher, err := sh.voucherRepos.GetByCode(ctx, i.VoucherCode)
			if err != nil {
				return nil, err
			}

			if i.StoreId != voucher.OwnerVoucher {
				return nil, responses.ErrUnableApplyVoucher
			}

			if err := sh.validateVoucherRequirement(ctx, req, voucher, i.SubTotal); err != nil {
				return nil, err
			}

			storeVoucher = append(storeVoucher, voucher)
		}

		if err := mapper.BindingStruct(storeVoucher, &resp.StoreVouchers); err != nil {
			return nil, err
		}
	}

	return &resp, nil
}

func (sh VoucherService) RollBackVoucher(ctx context.Context, req *dto.RollbackVoucherRequest) error {
	for _, i := range req.VoucherCodes {
		voucher, err := sh.voucherRepos.GetByCode(ctx, i)
		if err != nil {
			return err
		}

		voucher.VoucherCounts += 1

		if err := sh.voucherRepos.UpdateVoucherCounts(ctx, voucher); err != nil {
			return err
		}
	}

	return nil
}

func (sh VoucherService) CommitVoucherTransaction(ctx context.Context, msg *message.CreatePurchaseMessage) error {
	var msgReply []message.ReplyPurchaseMessage
	msgMap := make(map[string]int)
	var orderIds []string

	//// Initialize reply message
	for index, i := range msg.CheckoutData.OrderData {
		reply := message.ReplyPurchaseMessage{
			OrderID: i.OrderID,
			Status:  message.COMMIT_SUCCESS, // Assume success until an error occurs
		}

		orderIds = append(orderIds, i.OrderID)
		msgMap[i.StoreID] = index
		msgReply = append(msgReply, reply)
	}

	// Defer the reply message publishing
	defer func(replyPub *purchaseCreate.ReplyPurchaseTransactionPub, replyMsg []message.ReplyPurchaseMessage) {
		for _, i := range replyMsg {
			err := replyPub.ReplyPurchaseMessage(&i)
			if err != nil {
				log.Error(err)
			}
		}
	}(sh.voucherReply, msgReply)

	for _, i := range msg.Vouchers {
		var err error
		var orderIdLog []string

		voucher, err := sh.voucherRepos.GetByCode(ctx, i)
		if err != nil {
			log.Error(err)
		}

		if voucher.VoucherType == entities.STORE_DISCOUNT {
			orderIdLog = []string{msgReply[msgMap[voucher.OwnerVoucher]].OrderID}
		} else {
			orderIdLog = orderIds
		}

		err = sh.applyVoucherTransaction(ctx, voucher, orderIdLog, msg)

		if err != nil {
			switch voucher.VoucherType {
			case entities.STORE_DISCOUNT:
				msgReply[msgMap[voucher.OwnerVoucher]].Status = message.COMMIT_FAIL
			default:
				for index, _ := range msgReply {
					msgReply[index].Status = message.COMMIT_FAIL
				}
				break
			}
		}

	}

	return nil
}

func (sh VoucherService) applyVoucherTransaction(ctx context.Context, voucher *entities.Voucher, order []string, req *message.CreatePurchaseMessage) error {

	var err error

	if voucher.VoucherCounts < 0 {
		err = responses.ErrUnableApplyVoucher
	}

	usingLog := entities.VoucherUsingLog{
		VoucherCode: voucher.VoucherCode,
		VoucherID:   voucher.ID,
		UserID:      req.UserId,
		CheckoutPurchase: entities.CheckoutPurchase{
			CheckoutID: req.CheckoutData.CheckoutID,
			OrderIDs:   order,
		},
		Status: message.COMMIT_SUCCESS,
	}

	voucher.VoucherCounts-- //decrease voucher count

	if err = sh.voucherRepos.UpdateVoucherCounts(ctx, voucher); err != nil {
		return err
	}

	if err = sh.voucherRepos.CreateUsingVoucherLog(ctx, &usingLog); err != nil {
		return err
	}

	return nil
}

func (sh VoucherService) RollbackVoucherTransaction(ctx context.Context, req *message.RollbackPurchaseMessage) error {
	voucherLogs, err := sh.voucherRepos.FindVoucherLogByOrderID(ctx, req.OrderID)
	if err != nil {
		return err
	}

	for _, i := range voucherLogs {
		if len(i.CheckoutPurchase.OrderIDs) > 1 {
			voucherDetail, err := sh.voucherRepos.GetById(ctx, i.VoucherID.String())
			if err != nil {
				return err
			}

			voucherDetail.VoucherCounts--
			if err := sh.voucherRepos.UpdateVoucherCounts(ctx, voucherDetail); err != nil {
				return err
			}

			i.Status = message.COMMIT_FAIL
			if err := sh.voucherRepos.UpdateVoucherCounts(ctx, voucherDetail); err != nil {
				return err
			}
		}
	}

	return nil
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

func (sh VoucherService) UsingVoucherToOrder(ctx context.Context, req *dto.ApplyVoucherRequest) error {
	for _, i := range req.Vouchers {
		voucher, err := sh.voucherRepos.GetByCode(ctx, i)
		if err != nil {
			return err
		}

		if voucher.VoucherCounts < 0 {
			return responses.ErrUnableApplyVoucher
		}

		voucher.VoucherCounts -= 1
		if err := sh.voucherRepos.UpdateVoucherCounts(ctx, voucher); err != nil {
			return err
		}

		usingLog := entities.VoucherUsingLog{
			VoucherCode: voucher.VoucherCode,
			VoucherID:   voucher.ID,
			UserID:      req.UserId,
			CheckoutPurchase: entities.CheckoutPurchase{
				CheckoutID: req.CheckoutData.CheckoutID,
				OrderIDs:   req.CheckoutData.OrderIds,
			},
			Status: 0,
		}

		if err := sh.voucherRepos.CreateUsingVoucherLog(ctx, &usingLog); err != nil {
			return err
		}

	}

	return nil
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
