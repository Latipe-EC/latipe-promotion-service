package voucherserv

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"latipe-promotion-services/internal/domain/dto"
	"latipe-promotion-services/internal/domain/entities"
	responses "latipe-promotion-services/pkgs/response"
	"time"
)

func (sh VoucherService) validateVoucherRequirement(ctx context.Context, req *dto.PurchaseVoucherRequest,
	voucher *entities.Voucher, storeAmount ...int) error {
	//check time
	if voucher.VoucherCounts <= 0 || (!voucher.EndedTime.After(time.Now()) && !voucher.StatedTime.Before(time.Now()) ||
		voucher.Status != entities.ACTIVE) {
		return responses.ErrVoucherExpiredOrOutOfStock
	}

	//check required
	if voucher.VoucherRequire != nil {
		//check total
		if voucher.VoucherType != entities.STORE_DISCOUNT && int64(req.OrderTotalAmount) < voucher.VoucherRequire.MinRequire {
			return responses.ErrUnableApplyVoucher
		}

		if voucher.VoucherType != entities.STORE_DISCOUNT {
			if len(storeAmount) > 0 && int64(storeAmount[0]) < voucher.VoucherRequire.MinRequire {
				return responses.ErrUnableApplyVoucher
			}
		}

		//check payment method
		if voucher.VoucherRequire.PaymentMethod != 0 && req.PaymentMethod != voucher.VoucherRequire.PaymentMethod {
			return responses.ErrUnableApplyVoucher
		}

		//check number of using voucher
		total, err := sh.voucherRepos.CheckUsableVoucherOfUser(ctx, req.UserId, voucher.VoucherCode)
		if err != nil {
			return err
		}

		if total >= voucher.VoucherRequire.MaxVoucherPerUser {
			return responses.ErrUnableApplyVoucher
		}
	}
	return nil
}

func ParseStringToTime(dateStr string) time.Time {
	layout := "2006-01-02T15:04"
	// Get the local location
	location := time.Now().Location()

	date, err := time.ParseInLocation(layout, dateStr, location)
	if err != nil {
		return time.Now()
	}

	return date
}

func (sh VoucherService) validateVoucherRequest(ctx context.Context, req *dto.CreateVoucherRequest) error {
	currentTime := time.Now()
	startTime := ParseStringToTime(req.StatedTime)
	endedTime := ParseStringToTime(req.EndedTime)

	if req.VoucherCounts < 1 {
		return responses.ErrInvalidParameters
	}

	if req.VoucherType == entities.STORE_DISCOUNT && !sh.validateStoreVoucherPolicy(ctx, req) {
		return responses.ErrOutOfStorePolicy

	}

	if startTime.Before(currentTime) || endedTime.Before(currentTime) || startTime.After(endedTime) {
		return responses.ErrInvalidDatetime
	}

	//validate beween start time and ended time not over 30days
	OneMonthTime := 30 * 24 * time.Hour
	if endedTime.Sub(startTime) > OneMonthTime {
		return responses.ErrInvalidDatetime
	}

	if req.VoucherRequire.MinRequire < 0 {
		return responses.ErrBadRequest
	}

	switch req.VoucherType {
	case entities.FREE_SHIP:
		if req.DiscountData.ShippingValue == 0 {
			return responses.ErrInvalidVoucherData
		}
	default:
		if req.VoucherType == entities.PAYMENT_DISCOUNT && req.VoucherRequire.PaymentMethod == entities.COD_METHOD {
			return responses.ErrInvalidVoucherData
		}

		if req.DiscountData.DiscountType == entities.PERCENT_DISCOUNT &&
			(req.DiscountData.DiscountPercent == 0 || req.DiscountData.MaximumValue == 0) {
			return responses.ErrInvalidVoucherData
		}

		if req.DiscountData.DiscountType == entities.FIXED_DISCOUNT && (req.DiscountData.DiscountValue == 0) {
			return responses.ErrInvalidVoucherData
		}
	}

	return nil
}

func (sh VoucherService) validateStoreVoucherPolicy(ctx context.Context, req *dto.CreateVoucherRequest) bool {
	if req.VoucherCounts > entities.STORE_MAX_VC_COUNTS_POLICY {
		return false
	}

	if req.DiscountData.DiscountType == entities.FIXED_DISCOUNT &&
		req.DiscountData.DiscountValue > entities.STORE_MAX_FIXED_VALUE_POLICY {
		return false
	}

	if req.DiscountData.DiscountType == entities.PERCENT_DISCOUNT &&
		req.DiscountData.DiscountPercent > entities.STORE_MAX_PERCENT_POLICY &&
		req.DiscountData.MaximumValue > entities.STORE_MAX_FIXED_VALUE_POLICY {
		return false
	}

	totalCreatedCounts, err := sh.voucherRepos.CountAllVoucherCreatedInCurrentMonthByStoreId(ctx, req.StoreID)
	if err != nil {
		log.Error(err)
		return false
	}

	if totalCreatedCounts >= entities.STORE_MAX_CREATED_COUNTS_POLICY {
		return false
	}

	return true
}
