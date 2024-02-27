package voucherserv

import (
	"context"
	"latipe-promotion-services/internal/domain/dto"
	"latipe-promotion-services/internal/domain/entities"
	responses "latipe-promotion-services/pkgs/response"
	"time"
)

func (sh VoucherService) validateVoucherRequirement(ctx context.Context, req *dto.PurchaseVoucherRequest,
	voucher *entities.Voucher, storeAmount ...int) error {
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
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Now()
	}

	return date
}

func ValidateVoucherRequest(req *dto.CreateVoucherRequest) error {
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
