package voucherserv

import (
	"latipe-promotion-services/internal/domain/dto"
	"latipe-promotion-services/internal/domain/entities"
	responses "latipe-promotion-services/pkgs/response"
	"time"
)

func ParseStringToTime(dateStr string) time.Time {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Now()
	}

	return date
}

func ParseDateToString(date time.Time) string {
	layout := "2006-01-02"
	formattedTime := date.Format(layout)
	return formattedTime
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
