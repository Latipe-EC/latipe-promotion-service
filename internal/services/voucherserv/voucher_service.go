package voucherserv

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/mongo"
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

func (sh VoucherService) CheckoutPurchaseVoucherGrpc(ctx context.Context, req *dto.PurchaseVoucherRequest) (*dto.CheckoutPurchaseVoucherResponse, error) {
	resp := dto.CheckoutPurchaseVoucherResponse{}

	voucher, err := sh.voucherRepos.GetByCode(ctx, req.VoucherData.VoucherCode)
	if err != nil {
		return nil, err
	}

	if err := sh.validateVoucherRequirement(ctx, req, voucher); err != nil {
		return nil, err
	}

	if err := mapper.BindingStruct(voucher, &resp.VoucherDetail); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (sh VoucherService) CheckoutPurchase(ctx context.Context, req *dto.CheckoutVoucherRequest) (*dto.CheckoutVoucherResponse, error) {
	resp := dto.CheckoutVoucherResponse{}
	var voucherResp []dto.VoucherUserDetail

	for _, i := range req.Vouchers {
		var respData dto.VoucherUserDetail
		voucher, err := sh.voucherRepos.GetByCode(ctx, i)
		if err != nil {
			return nil, err
		}

		if err := mapper.BindingStruct(voucher, &respData); err != nil {
			return nil, err
		}

		totalVouchers, err := sh.voucherRepos.CheckUsableVoucherOfUser(ctx, req.UserID, voucher.VoucherCode)
		if err != nil {
			return nil, err
		}

		if totalVouchers >= voucher.VoucherRequire.MaxVoucherPerUser {
			respData.Usable = false
		} else {
			respData.Usable = true
		}

		respData.CountUsable = voucher.VoucherRequire.MaxVoucherPerUser - totalVouchers
		voucherResp = append(voucherResp, respData)
	}
	resp.Items = voucherResp

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
	// Initialize reply message
	reply := message.ReplyPurchaseMessage{
		OrderID: msg.OrderID,
		Status:  message.COMMIT_FAIL, // Assume success until an error occurs
	}
	defer func(replyPub *purchaseCreate.ReplyPurchaseTransactionPub) {
		err := replyPub.ReplyPurchaseMessage(&reply)
		if err != nil {
			log.Error(err)
		}
	}(sh.voucherReply)

	for _, code := range msg.VoucherCodes {
		voucher, err := sh.voucherRepos.GetByCode(ctx, code)
		if err != nil {
			log.Error(err)
			return err
		}

		//init using voucher log
		usingLog := entities.VoucherUsingLog{
			VoucherCode: code,
			VoucherID:   voucher.ID,
			UserID:      msg.UserID,
			CheckoutPurchase: entities.CheckoutPurchase{
				CheckoutID: msg.CheckoutID,
				OrderIDs:   []string{msg.OrderID},
			},
			Status: message.COMMIT_SUCCESS,
		}

		switch voucher.VoucherType {
		case entities.FREE_SHIP, entities.STORE_DISCOUNT:
			err = sh.applyVoucherTransaction(ctx, voucher, &usingLog)
			if err != nil {
				log.Error(err)
			}

		case entities.PAYMENT_DISCOUNT:
			existLog, err := sh.voucherRepos.FindVoucherLogByVoucherCodeAndCheckoutID(ctx, code, msg.CheckoutID)
			if err != nil {
				log.Error(err)
				switch {
				case errors.Is(err, mongo.ErrNoDocuments):
					err = sh.applyVoucherTransaction(ctx, voucher, &usingLog)
					if err != nil {
						log.Error(err)
					}
				default:
					return err
				}
			} else {
				existLog.CheckoutPurchase.OrderIDs = append(existLog.CheckoutPurchase.OrderIDs, msg.OrderID)
				if err = sh.updatePaymentVoucherTransaction(ctx, voucher, existLog); err != nil {
					return err
				}
			}
		}
	}

	reply.Status = message.COMMIT_SUCCESS
	return nil
}

func (sh VoucherService) applyVoucherTransaction(ctx context.Context, voucher *entities.Voucher, voucherLog *entities.VoucherUsingLog) error {

	var err error

	if voucher.VoucherCounts < 0 {
		err = responses.ErrUnableApplyVoucher
	}
	voucher.VoucherCounts-- //decrease voucher count

	if err = sh.voucherRepos.UpdateVoucherCounts(ctx, voucher); err != nil {
		return err
	}

	if err = sh.voucherRepos.CreateUsingVoucherLog(ctx, voucherLog); err != nil {
		return err
	}

	return nil
}

func (sh VoucherService) updatePaymentVoucherTransaction(ctx context.Context, voucher *entities.Voucher, voucherLog *entities.VoucherUsingLog) error {

	var err error

	if voucher.VoucherCounts < 0 {
		err = responses.ErrUnableApplyVoucher
	}
	voucher.VoucherCounts-- //decrease voucher count

	if err = sh.voucherRepos.UpdateVoucherCounts(ctx, voucher); err != nil {
		return err
	}

	if err = sh.voucherRepos.UpdateUsingVoucherLog(ctx, voucherLog); err != nil {
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

func (sh VoucherService) GetComingVoucher(ctx context.Context, query *pagable.Query) (*pagable.ListResponse, error) {
	vouchers, total, err := sh.voucherRepos.GetComingVoucher(ctx, query)
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
		totalVouchers, err := sh.voucherRepos.CheckUsableVoucherOfUser(ctx, userId, i.VoucherCode)
		if err != nil {
			return nil, err
		}

		if totalVouchers >= i.VoucherRequire.MaxVoucherPerUser {
			voucherResp[index].Usable = false
		} else {
			voucherResp[index].Usable = true
		}

		voucherResp[index].CountUsable = i.VoucherRequire.MaxVoucherPerUser - totalVouchers
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
